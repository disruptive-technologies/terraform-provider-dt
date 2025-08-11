// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &contactGroupResource{}
	_ resource.ResourceWithConfigure   = &contactGroupResource{}
	_ resource.ResourceWithImportState = &contactGroupResource{}
)

// NewContactGroupResource creates a new resource for managing contact groups.
func NewContactGroupResource() resource.Resource {
	return &contactGroupResource{}
}

// contactGroupResource is a Terraform provider for managing contact groups.
type contactGroupResource struct {
	client *dt.Client
}

// Metadata returns the resource type name

func (r *contactGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact_group"
}

// ImportState imports the state of a contact group resource.
func (r *contactGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Schema defines the schema for the resource.
func (r *contactGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed: true,
				Description: `The resource name of the contact group.
								Format is "organizations/{organization}/contactGroups/{contact_group}".`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "The organization ID of the contact group.",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the contact group.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the contact group.",
				Default:     stringdefault.StaticString(""),
			},
			"contact_count": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of contacts in the group.",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type contactGroupResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Organization types.String `tfsdk:"organization"`
	DisplayName  types.String `tfsdk:"display_name"`
	Description  types.String `tfsdk:"description"`
	ContactCount types.Int32  `tfsdk:"contact_count"`
}

// Create creates the resource and sets the initial state.
func (r *contactGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read the plan to get the contact group details
	var plan contactGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the plan to a CreateContactGroupRequest
	createRequest, diags := stateToCreateContactGroupRequest(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the contact group using the api client
	createdGroup, err := r.client.CreateContactGroup(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create contact group",
			"An error occurred while creating the contact group: "+err.Error(),
		)
		return
	}

	// Convert the created contact group to the resource model
	newState, diags := contactGroupToState(createdGroup)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Set the state with the created contact group
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
func (r *contactGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read the current state of the resource
	var state contactGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the contact group from the API
	contactGroup, err := r.client.GetContactGroup(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read contact group",
			"An error occurred while reading the contact group: "+err.Error(),
		)
		return
	}

	// Convert the contact group to the resource model
	newState, diags := contactGroupToState(contactGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state with the updated contact group
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource with the provided plan.
// It compares the plan with the current state and applies changes accordingly.
func (r *contactGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read the current state of the resource
	var plan contactGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the plan to an UpdateContactGroupRequest
	updateRequest, diags := stateToUpdateContactGroupRequest(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the contact group using the API client
	updatedGroup, err := r.client.UpdateContactGroup(ctx, updateRequest, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update contact group",
			"An error occurred while updating the contact group: "+err.Error(),
		)
		return
	}

	// Convert the updated contact group to the resource model
	state, diags := contactGroupToState(updatedGroup)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set the state with the updated contact group
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *contactGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Read the current state of the resource
	var state contactGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the contact group using the API client
	err := r.client.DeleteContactGroup(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete contact group",
			"An error occurred while deleting the contact group: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *contactGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func contactGroupToState(contactGroup dt.ContactGroup) (contactGroupResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	organizationID, _, err := dt.ParseResourceName(contactGroup.Name)
	if err != nil {
		diags = append(diags, diag.NewAttributeErrorDiagnostic(
			path.Root("name"),
			"Failed to parse contact group name",
			"An error occurred while parsing the contact group name: "+err.Error(),
		))
	}

	return contactGroupResourceModel{
		Name:         types.StringValue(contactGroup.Name),
		Organization: types.StringValue("organizations/" + organizationID),
		DisplayName:  types.StringValue(contactGroup.DisplayName),
		Description:  types.StringValue(contactGroup.Description),
		ContactCount: types.Int32Value(contactGroup.ContactCount),
	}, diags
}

func stateToCreateContactGroupRequest(state contactGroupResourceModel) (dt.CreateContactGroupRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	if state.Organization.IsNull() || state.Organization.IsUnknown() {
		diags = append(diags, diag.NewAttributeErrorDiagnostic(
			path.Root("organization"),
			"Organization is required",
			"An error occurred while processing the contact group: organization is required.",
		))
	}

	return dt.CreateContactGroupRequest{
		Organization: state.Organization.ValueString(),
		ContactGroup: dt.ContactGroup{
			DisplayName: state.DisplayName.ValueString(),
			Description: state.Description.ValueString(),
		},
	}, diags
}

func stateToUpdateContactGroupRequest(state contactGroupResourceModel) (dt.UpdateContactGroupRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	if state.DisplayName.IsNull() || state.DisplayName.IsUnknown() {
		diags = append(diags, diag.NewAttributeErrorDiagnostic(
			path.Root("display_name"),
			"Display name is required",
			"An error occurred while processing the contact group: display name is required.",
		))
	}

	return dt.UpdateContactGroupRequest{
		DisplayName: state.DisplayName.ValueString(),
		Description: state.Description.ValueString(),
	}, diags
}
