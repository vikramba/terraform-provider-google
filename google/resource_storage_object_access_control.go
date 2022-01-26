// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStorageObjectAccessControl() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageObjectAccessControlCreate,
		Read:   resourceStorageObjectAccessControlRead,
		Update: resourceStorageObjectAccessControlUpdate,
		Delete: resourceStorageObjectAccessControlDelete,

		Importer: &schema.ResourceImporter{
			State: resourceStorageObjectAccessControlImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The name of the bucket.`,
			},
			"entity": {
				Type:     schema.TypeString,
				Required: true,
				Description: `The entity holding the permission, in one of the following forms:
  * user-{{userId}}
  * user-{{email}} (such as "user-liz@example.com")
  * group-{{groupId}}
  * group-{{email}} (such as "group-example@googlegroups.com")
  * domain-{{domain}} (such as "domain-example.com")
  * project-team-{{projectId}}
  * allUsers
  * allAuthenticatedUsers`,
			},
			"object": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the object to apply the access control to.`,
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateEnum([]string{"OWNER", "READER"}),
				Description:  `The access permission for the entity. Possible values: ["OWNER", "READER"]`,
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain associated with the entity.`,
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The email address associated with the entity.`,
			},
			"entity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID for the entity`,
			},
			"generation": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The content generation of the object, if applied to an object.`,
			},
			"project_team": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The project team associated with the entity`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The project team associated with the entity`,
						},
						"team": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateEnum([]string{"editors", "owners", "viewers", ""}),
							Description:  `The team. Possible values: ["editors", "owners", "viewers"]`,
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func resourceStorageObjectAccessControlCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	bucketProp, err := expandStorageObjectAccessControlBucket(d.Get("bucket"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("bucket"); !isEmptyValue(reflect.ValueOf(bucketProp)) && (ok || !reflect.DeepEqual(v, bucketProp)) {
		obj["bucket"] = bucketProp
	}
	entityProp, err := expandStorageObjectAccessControlEntity(d.Get("entity"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("entity"); !isEmptyValue(reflect.ValueOf(entityProp)) && (ok || !reflect.DeepEqual(v, entityProp)) {
		obj["entity"] = entityProp
	}
	objectProp, err := expandStorageObjectAccessControlObject(d.Get("object"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("object"); !isEmptyValue(reflect.ValueOf(objectProp)) && (ok || !reflect.DeepEqual(v, objectProp)) {
		obj["object"] = objectProp
	}
	roleProp, err := expandStorageObjectAccessControlRole(d.Get("role"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("role"); !isEmptyValue(reflect.ValueOf(roleProp)) && (ok || !reflect.DeepEqual(v, roleProp)) {
		obj["role"] = roleProp
	}

	lockName, err := replaceVars(d, config, "storage/buckets/{{bucket}}/objects/{{object}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{StorageBasePath}}b/{{bucket}}/o/{{%object}}/acl")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new ObjectAccessControl: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating ObjectAccessControl: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{bucket}}/{{object}}/{{entity}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating ObjectAccessControl %q: %#v", d.Id(), res)

	return resourceStorageObjectAccessControlRead(d, meta)
}

func resourceStorageObjectAccessControlRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{StorageBasePath}}b/{{bucket}}/o/{{%object}}/acl/{{entity}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("StorageObjectAccessControl %q", d.Id()))
	}

	if err := d.Set("bucket", flattenStorageObjectAccessControlBucket(res["bucket"], d, config)); err != nil {
		return fmt.Errorf("Error reading ObjectAccessControl: %s", err)
	}
	if err := d.Set("domain", flattenStorageObjectAccessControlDomain(res["domain"], d, config)); err != nil {
		return fmt.Errorf("Error reading ObjectAccessControl: %s", err)
	}
	if err := d.Set("email", flattenStorageObjectAccessControlEmail(res["email"], d, config)); err != nil {
		return fmt.Errorf("Error reading ObjectAccessControl: %s", err)
	}
	if err := d.Set("entity", flattenStorageObjectAccessControlEntity(res["entity"], d, config)); err != nil {
		return fmt.Errorf("Error reading ObjectAccessControl: %s", err)
	}
	if err := d.Set("entity_id", flattenStorageObjectAccessControlEntityId(res["entityId"], d, config)); err != nil {
		return fmt.Errorf("Error reading ObjectAccessControl: %s", err)
	}
	if err := d.Set("generation", flattenStorageObjectAccessControlGeneration(res["generation"], d, config)); err != nil {
		return fmt.Errorf("Error reading ObjectAccessControl: %s", err)
	}
	if err := d.Set("object", flattenStorageObjectAccessControlObject(res["object"], d, config)); err != nil {
		return fmt.Errorf("Error reading ObjectAccessControl: %s", err)
	}
	if err := d.Set("project_team", flattenStorageObjectAccessControlProjectTeam(res["projectTeam"], d, config)); err != nil {
		return fmt.Errorf("Error reading ObjectAccessControl: %s", err)
	}
	if err := d.Set("role", flattenStorageObjectAccessControlRole(res["role"], d, config)); err != nil {
		return fmt.Errorf("Error reading ObjectAccessControl: %s", err)
	}

	return nil
}

func resourceStorageObjectAccessControlUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	obj := make(map[string]interface{})
	bucketProp, err := expandStorageObjectAccessControlBucket(d.Get("bucket"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("bucket"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, bucketProp)) {
		obj["bucket"] = bucketProp
	}
	entityProp, err := expandStorageObjectAccessControlEntity(d.Get("entity"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("entity"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, entityProp)) {
		obj["entity"] = entityProp
	}
	objectProp, err := expandStorageObjectAccessControlObject(d.Get("object"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("object"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, objectProp)) {
		obj["object"] = objectProp
	}
	roleProp, err := expandStorageObjectAccessControlRole(d.Get("role"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("role"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, roleProp)) {
		obj["role"] = roleProp
	}

	lockName, err := replaceVars(d, config, "storage/buckets/{{bucket}}/objects/{{object}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{StorageBasePath}}b/{{bucket}}/o/{{%object}}/acl/{{entity}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating ObjectAccessControl %q: %#v", d.Id(), obj)

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "PUT", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating ObjectAccessControl %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating ObjectAccessControl %q: %#v", d.Id(), res)
	}

	return resourceStorageObjectAccessControlRead(d, meta)
}

func resourceStorageObjectAccessControlDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	lockName, err := replaceVars(d, config, "storage/buckets/{{bucket}}/objects/{{object}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{StorageBasePath}}b/{{bucket}}/o/{{%object}}/acl/{{entity}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting ObjectAccessControl %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "ObjectAccessControl")
	}

	log.Printf("[DEBUG] Finished deleting ObjectAccessControl %q: %#v", d.Id(), res)
	return nil
}

func resourceStorageObjectAccessControlImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"(?P<bucket>[^/]+)/(?P<object>.+)/(?P<entity>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{bucket}}/{{object}}/{{entity}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenStorageObjectAccessControlBucket(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenStorageObjectAccessControlDomain(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageObjectAccessControlEmail(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageObjectAccessControlEntity(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageObjectAccessControlEntityId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageObjectAccessControlGeneration(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := stringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenStorageObjectAccessControlObject(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageObjectAccessControlProjectTeam(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["project_number"] =
		flattenStorageObjectAccessControlProjectTeamProjectNumber(original["projectNumber"], d, config)
	transformed["team"] =
		flattenStorageObjectAccessControlProjectTeamTeam(original["team"], d, config)
	return []interface{}{transformed}
}
func flattenStorageObjectAccessControlProjectTeamProjectNumber(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageObjectAccessControlProjectTeamTeam(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageObjectAccessControlRole(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandStorageObjectAccessControlBucket(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandStorageObjectAccessControlEntity(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandStorageObjectAccessControlObject(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandStorageObjectAccessControlRole(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
