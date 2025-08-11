// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &contactResource{}
	_ resource.ResourceWithConfigure   = &contactResource{}
	_ resource.ResourceWithImportState = &contactResource{}
)

// NewContactResource creates a new contact resource.
func NewContactResource() resource.Resource {
	return &contactResource{}
}

// contactResource is a Terraform resource for managing contacts.
type contactResource struct {
	client *dt.Client
}

// Metadata returns the resource type name
func (r *contactResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact"
}

// ImportState imports the contact resource state.
func (r *contactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Schema defines the schema for the resource
func (r *contactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: `The resource name of the contact. Format is "projects/{project}/contacts/{contact}".`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"contact_group": schema.StringAttribute{
				Required: true,
				Description: `The resource name of the contact group this contact belongs to.
  								Format is "organizations/{organization}/contactGroups/{contact_group}".`,
			},
			"project": schema.StringAttribute{
				Required:    true,
				Description: `The resource name of the project this contact belongs to. Format is "projects/{project}".`,
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: `The display name of the contact.`,
				Validators:  []validator.String{stringvalidator.LengthBetween(1, 100)},
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: `The email address of the contact.`,
				Default:     stringdefault.StaticString(""),
			},
			"phone_number": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: `The phone number of the contact.`,
				Default:     stringdefault.StaticString(""),
			},
			"has_project_access": schema.BoolAttribute{
				Computed:    true,
				Description: `Indicates whether the contact has access to the project.`,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type contactResourceModel struct {
	Name             types.String `tfsdk:"name"`
	ContactGroup     types.String `tfsdk:"contact_group"`
	Project          types.String `tfsdk:"project"`
	DisplayName      types.String `tfsdk:"display_name"`
	Email            types.String `tfsdk:"email"`
	PhoneNumber      types.String `tfsdk:"phone_number"`
	HasProjectAccess types.Bool   `tfsdk:"has_project_access"`
}

// Create creates the resource and sets the initial state.
func (r *contactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read the plan to get the contact details.
	var plan contactResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the plan to a CreateContactRequest.
	createRequest := stateToCreateContactRequest(plan)

	// Create the contact using the client.
	createdContact, err := r.client.CreateContact(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create contact",
			"An error occurred while creating the contact: "+err.Error(),
		)
		return
	}

	// Convert the created contact to the resource model.
	createdModel, diags := contactToState(createdContact)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Set the resource state with the created contact.
	diags = resp.State.Set(ctx, createdModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *contactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read the current state of the resource.
	var state contactResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the contact using the client.
	contact, err := r.client.GetContact(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read contact",
			"An error occurred while reading the contact: "+err.Error(),
		)
		return
	}

	// Convert the contact to the resource model.
	state, diags = contactToState(contact)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Set the resource state with the updated contact.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource with the provided plan.
func (r *contactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read the plan to get the updated contact details.
	var plan contactResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the plan to an UpdateContactRequest.
	updateRequest := stateToUpdateContactRequest(plan)

	// Update the contact using the client.
	updatedContact, err := r.client.UpdateContact(ctx, updateRequest, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update contact",
			"An error occurred while updating the contact: "+err.Error(),
		)
		return
	}

	// Convert the updated contact to the resource model.
	updatedModel, diags := contactToState(updatedContact)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Set the resource state with the updated contact.
	diags = resp.State.Set(ctx, updatedModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource.
func (r *contactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Read the current state of the resource.
	var state contactResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the contact using the client.
	err := r.client.DeleteContact(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete contact",
			"An error occurred while deleting the contact: "+err.Error(),
		)
		return
	}
}

func (r *contactResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dt.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *dt.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// stateToContact converts the resource model to the dt api Contact type.
func stateToContact(model contactResourceModel) dt.Contact {
	return dt.Contact{
		Name:             model.Name.ValueString(),
		ContactGroup:     model.ContactGroup.ValueString(),
		DisplayName:      model.DisplayName.ValueString(),
		Email:            model.Email.ValueString(),
		PhoneNumber:      model.PhoneNumber.ValueString(),
		HasProjectAccess: model.HasProjectAccess.ValueBool(),
	}
}

func contactToState(contact dt.Contact) (contactResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	projectID, _, err := dt.ParseResourceName(contact.Name)
	if err != nil {
		diags = append(diags, diag.NewErrorDiagnostic(
			"Failed to parse contact name",
			"An error occurred while parsing the contact name: "+err.Error(),
		))
	}
	return contactResourceModel{
		Name:             types.StringValue(contact.Name),
		Project:          types.StringValue("projects/" + projectID),
		ContactGroup:     types.StringValue(contact.ContactGroup),
		DisplayName:      types.StringValue(contact.DisplayName),
		Email:            types.StringValue(contact.Email),
		PhoneNumber:      types.StringValue(contact.PhoneNumber),
		HasProjectAccess: types.BoolValue(contact.HasProjectAccess),
	}, diags
}

func stateToCreateContactRequest(model contactResourceModel) dt.CreateContactRequest {
	return dt.CreateContactRequest{
		Project: model.Project.ValueString(),
		Contact: stateToContact(model),
	}
}

func stateToUpdateContactRequest(model contactResourceModel) dt.UpdateContactRequest {
	return dt.UpdateContactRequest{
		ContactGroup: model.ContactGroup.ValueString(),
		DisplayName:  model.DisplayName.ValueString(),
		Email:        model.Email.ValueString(),
		PhoneNumber:  model.PhoneNumber.ValueString(),
	}
}
