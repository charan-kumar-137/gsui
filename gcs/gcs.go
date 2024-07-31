package gcs

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/api/iterator"
)

var (
	fieldStyle = lipgloss.NewStyle().Bold(true)
	valueStyle = lipgloss.NewStyle()
	linkStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#3367D6")).Underline(true).Italic(true)

	bucketCols = []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Created", Width: 20},
		{Title: "LocationType", Width: 15},
		{Title: "Location", Width: 10},
		{Title: "DefaultStorageClass", Width: 20},
	}

	objectCols = []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Size", Width: 10},
		{Title: "Type", Width: 10},
		{Title: "Storage Class", Width: 15},
		{Title: "Created", Width: 20},
		{Title: "Last Modified", Width: 20},
	}

	projectId string
	ctx       context.Context
	client    *storage.Client

	// buckets []*Bucket = []*Bucket{
	// 	{name: "test1", created: "18:00", locationType: "Multi Region", location: "US", defaultStorageClass: "Standard"},
	// 	{name: "test2", created: "19:01", locationType: "Multi Region", location: "AU", defaultStorageClass: "Standard"},
	// 	{name: "test3", created: "16:02", locationType: "Multi Region", location: "EU", defaultStorageClass: "Near Line"},
	// 	{name: "test4", created: "14:03", locationType: "Multi Region", location: "APAC", defaultStorageClass: "Standard"},
	// }

	// objects []*Object = []*Object{
	// 	{name: "O1", size: "100", objectType: "go", created: "18:00", lastModified: "18:01", storageClass: "Standard", customTime: "a", objectRetainUntil: "adasdaaaaaaaaaaaa"},
	// 	{name: "O2", size: "100", objectType: "txt", created: "18:00", lastModified: "18:01", storageClass: "Standard"},
	// 	{name: "O3", size: "100", objectType: "py", created: "18:00", lastModified: "18:01", storageClass: "Standard"},
	// 	{name: "O4", size: "", objectType: "", created: "18:00", lastModified: "18:01", storageClass: "Standard"},
	// 	{name: "O5", size: "100", objectType: "go", created: "18:00", lastModified: "18:01", storageClass: "Standard"},
	// }

	// bucketObjectMap map[*Bucket][]*Object = map[*Bucket][]*Object{
	// 	buckets[0]: {
	// 		objects[0],
	// 		objects[1],
	// 	},
	// 	buckets[1]: {
	// 		objects[2],
	// 	},
	// 	buckets[2]: {
	// 		objects[3],
	// 	},
	// 	buckets[3]: {
	// 		objects[4],
	// 	},
	// }
)

func Init(_projectId string) error {
	projectId = _projectId
	var err error
	ctx = context.Background()
	client, err = storage.NewClient(ctx)
	return err
}

type Bucket struct {
	name                string
	created             string
	locationType        string
	location            string
	defaultStorageClass string
	lastModified        string
	publicAccess        string
	accessControl       string
	protection          string
	bucketRetention     string
	lifeCycleRules      string
	tags                string
	encryption          string
	labels              string
	requesterPays       string
	replication         string
}

func (b Bucket) GetName() string {
	return b.name
}

func renderFieldValue(field, value string) string {
	return lipgloss.NewStyle().Render(fieldStyle.Render(field) + " " + valueStyle.Render(value) + "\n")
}

func renderFieldHref(field, href string) string {
	return lipgloss.NewStyle().Render(fieldStyle.Render(field) + " " + linkStyle.Render(href) + "\n")
}

func renderIndent(indent int, str string) string {
	return strings.Repeat("\t", indent) + str
}

func (b Bucket) DisplayString() string {
	var sb strings.Builder

	sb.WriteString(renderFieldValue("Last Modified:", b.lastModified))
	sb.WriteString(renderFieldValue("Public Access:", b.publicAccess))
	sb.WriteString(renderFieldValue("Access Control:", b.accessControl))
	sb.WriteString(renderFieldValue("Protection:", b.protection))
	sb.WriteString(renderFieldValue("Bucket Retention:", b.bucketRetention))
	sb.WriteString(renderFieldValue("Life Cycle Rules:", b.lifeCycleRules))
	sb.WriteString(renderFieldValue("Tags:", b.tags))
	sb.WriteString(renderFieldValue("Encryption:", b.encryption))
	sb.WriteString(renderFieldValue("Labels:", b.labels))
	sb.WriteString(renderFieldValue("Requester Pays:", b.requesterPays))
	sb.WriteString(renderFieldValue("Replication:", b.replication))

	return sb.String()
}

type Object struct {
	name              string
	size              string
	objectType        string
	created           string
	lastModified      string
	storageClass      string
	customTime        string
	publicURL         string
	authenticatedURL  string
	gsutilURI         string
	publicAccess      string
	versionHistory    string
	objectRetainUntil string
	bucketRetainUntil string
	holdStatus        string
	encryptionType    string
}

func (o Object) GetName() string {
	return o.name
}

func (o Object) DisplayString() string {
	var sb strings.Builder

	var currentIndent = 0
	sb.WriteString(renderFieldValue("Overview", ""))
	currentIndent += 1
	sb.WriteString(renderIndent(currentIndent, renderFieldValue("Custom Time:", o.customTime)))
	sb.WriteString(renderIndent(currentIndent, renderFieldHref("Public URL:", o.publicURL)))
	sb.WriteString(renderIndent(currentIndent, renderFieldHref("Authenticated URL:", o.authenticatedURL)))
	sb.WriteString(renderIndent(currentIndent, renderFieldHref("gsutil URI:", o.gsutilURI)))
	currentIndent -= 1

	sb.WriteString(renderFieldValue("Permissions", ""))

	sb.WriteString(renderIndent(currentIndent+1, renderFieldValue("Public Access:", o.publicAccess)))

	sb.WriteString(renderFieldValue("Protection", ""))
	currentIndent += 1
	sb.WriteString(renderIndent(currentIndent, renderFieldValue("Version History:", o.versionHistory)))
	sb.WriteString(renderIndent(currentIndent, renderFieldValue("Retention Expiration Time", "")))
	sb.WriteString(renderIndent(currentIndent+1, renderFieldValue("Object Retention retain until time:", o.objectRetainUntil)))
	sb.WriteString(renderIndent(currentIndent+1, renderFieldValue("Bucket Retention retain until time:", o.bucketRetainUntil)))
	sb.WriteString(renderIndent(currentIndent, renderFieldValue("Hold Status:", o.holdStatus)))
	sb.WriteString(renderIndent(currentIndent, renderFieldValue("Encryption Type:", o.encryptionType)))
	currentIndent -= 1

	return sb.String()
}

type DataError struct {
	msg string
}

func (de DataError) Error() string {
	return de.msg
}

type Data struct {
	IsBucket bool
	buckets  []*Bucket
	objects  []*Object
	err      error
}

func (data Data) GetError() error {
	return data.err
}

func (data Data) GetBucket(index int) *Bucket {

	if index < 0 || index >= len(data.buckets) {
		return nil
	}

	return data.buckets[index]
}

func (data Data) GetObject(index int) *Object {
	if index < 0 || index >= len(data.objects) {
		return nil
	}

	return data.objects[index]
}

func (data Data) GetTableData() ([]table.Column, []table.Row) {
	if data.IsBucket {
		return bucketCols, convertBucketToRows(data.buckets)
	} else {
		return objectCols, convertObjectToRows(data.objects)
	}
}

func convertBucketToRows(buckets []*Bucket) []table.Row {
	var rows []table.Row

	for _, bucket := range buckets {
		rows = append(rows, table.Row{bucket.name, bucket.created, bucket.locationType, bucket.location, bucket.defaultStorageClass})
	}

	return rows
}

func convertObjectToRows(objects []*Object) []table.Row {
	var rows []table.Row

	for _, object := range objects {
		rows = append(rows, table.Row{object.name, object.size, object.objectType, object.storageClass, object.created, object.lastModified})
	}

	return rows
}

func listToString(list []string, sep string) string {
	var sb strings.Builder

	for i, item := range list {
		sb.WriteString(item)
		if i != len(list)-1 {
			sb.WriteString(sep)
		}
	}

	return sb.String()
}

func GetData(path string) *Data {
	if len(path) == 0 {
		// Get All the Buckets in Project
		var buckets []*Bucket = make([]*Bucket, 0)
		it := client.Buckets(ctx, projectId)
		for {
			bucketAttrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				// return err TODO
			} else {
				bucket := getBucket(bucketAttrs)
				buckets = append(buckets, bucket)
			}
		}

		return &Data{IsBucket: true, buckets: buckets}
	}

	var bucket string

	var index = strings.Index(path, "/")

	if index == -1 {
		bucket = path
	} else {
		bucket = path[:index]
	}

	_, err := client.Bucket(bucket).Attrs(ctx)

	if err == storage.ErrBucketNotExist {
		// return &Data{err: DataError{"Bucket Not Found - " + bucket}}
		return nil
	}

	if err != nil {
		// Handle Error TODO
		return nil
	}

	var prefix string = ""

	if index == -1 {
		prefix = ""
	} else {
		prefix = path[index+1:]
	}

	query := &storage.Query{Prefix: prefix}
	var objects []*Object

	it := client.Bucket(bucket).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: log.Fatal(err)
		} else {
			object := newFunction(attrs)
			objects = append(objects, object)
		}
	}

	if len(objects) != 0 {
		return &Data{IsBucket: false, objects: objects}
	}

	return nil
}

func newFunction(attrs *storage.ObjectAttrs) *Object {
	var authenticatedURL string = fmt.Sprintf("https://storage.cloud.google.com/%s/%s", attrs.Bucket, attrs.Name)
	var gsutilURI string = fmt.Sprintf("gs://%s/%s", attrs.Bucket, attrs.Name)
	var encryption string = "Google Managed"

	if len(attrs.CustomerKeySHA256) != 0 {
		encryption = attrs.CustomerKeySHA256
	} else if len(attrs.KMSKeyName) != 0 {
		encryption = attrs.KMSKeyName
	}

	object := &Object{
		name:             attrs.Name,
		size:             fmt.Sprint(attrs.Size),
		objectType:       attrs.ContentType,
		created:          attrs.Created.String(),
		lastModified:     attrs.Updated.String(),
		storageClass:     attrs.StorageClass,
		customTime:       attrs.CustomTime.String(),
		publicURL:        "NA",
		authenticatedURL: authenticatedURL,
		gsutilURI:        gsutilURI,
		// publicAccess      string
		// versionHistory    string
		// objectRetainUntil string
		// bucketRetainUntil string
		// holdStatus        string
		encryptionType: encryption,
	}
	return object
}

func getBucket(bucketAttrs *storage.BucketAttrs) *Bucket {
	var accessControl string = "Uniform"
	var protection string = "None"
	var bucketRetention = "None"
	var lifeCycleRules string = "None"
	var tags []string
	var labels []string
	var requesterPays = "OFF"
	var encryption = "Google Managed"

	if !bucketAttrs.UniformBucketLevelAccess.Enabled {
		accessControl = "Fine Grained"
	}

	if bucketAttrs.VersioningEnabled {
		protection = "Versioning"
	}

	if len(bucketAttrs.Lifecycle.Rules) != 0 {
		lifeCycleRules = fmt.Sprint(len(bucketAttrs.Lifecycle.Rules), " Rules")
	}

	for k, v := range bucketAttrs.Labels {
		labels = append(labels, k+": "+v)
	}

	if bucketAttrs.RequesterPays {
		requesterPays = "ON"
	}

	if bucketAttrs.Encryption != nil {
		encryption = bucketAttrs.Encryption.DefaultKMSKeyName
	}

	bucket := &Bucket{
		name:                bucketAttrs.Name,
		created:             bucketAttrs.Created.String(),
		locationType:        bucketAttrs.LocationType,
		location:            bucketAttrs.Location,
		defaultStorageClass: bucketAttrs.StorageClass,
		lastModified:        "NA",
		publicAccess:        bucketAttrs.PublicAccessPrevention.String(),
		accessControl:       accessControl,
		protection:          protection,
		bucketRetention:     bucketRetention,
		lifeCycleRules:      lifeCycleRules,
		tags:                listToString(tags, ","),
		encryption:          encryption,
		labels:              listToString(labels, ","),
		requesterPays:       requesterPays,
		replication:         bucketAttrs.RPO.String(),
	}
	return bucket
}
