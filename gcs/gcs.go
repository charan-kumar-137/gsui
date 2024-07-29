package gcs

import "github.com/charmbracelet/bubbles/table"

var (
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

	buckets []*Bucket = []*Bucket{
		{"test1", "18:00", "Multi Region", "US", "Standard"},
		{"test2", "19:01", "Multi Region", "AU", "Standard"},
		{"test3", "16:02", "Multi Region", "EU", "Near Line"},
		{"test4", "14:03", "Multi Region", "APAC", "Standard"},
	}

	objects []*Object = []*Object{
		{"O1", "100", "go", "18:00", "18:01", "Standard"},
		{"O2", "100", "txt", "18:00", "18:01", "Standard"},
		{"O3", "100", "py", "18:00", "18:01", "Standard"},
		{"O4", "", "", "18:00", "18:01", "Standard"},
		{"O5", "100", "go", "18:00", "18:01", "Standard"},
	}

	bucketObjectMap map[*Bucket][]*Object = map[*Bucket][]*Object{
		buckets[0]: {
			objects[0],
			objects[1],
		},
		buckets[1]: {
			objects[2],
		},
		buckets[2]: {
			objects[3],
		},
		buckets[3]: {
			objects[4],
		},
	}
)

type Bucket struct {
	Name                string
	Created             string
	LoacationType       string
	Location            string
	DefaultStorageClass string
}

type Object struct {
	Name         string
	Size         string
	Type         string
	Created      string
	LastModified string
	StorageClass string
}

type Data struct {
	IsBucket bool
	buckets  []*Bucket
	objects  []*Object
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
		rows = append(rows, table.Row{bucket.Name, bucket.Created, bucket.LoacationType, bucket.Location, bucket.DefaultStorageClass})
	}

	return rows
}

func convertObjectToRows(objects []*Object) []table.Row {
	var rows []table.Row

	for _, object := range objects {
		rows = append(rows, table.Row{object.Name, object.Size, object.Type, object.StorageClass, object.Created, object.LastModified})
	}

	return rows
}

func GetData(path string) *Data {
	if len(path) == 0 {
		return &Data{IsBucket: true, buckets: buckets}
	}

	for _, bucket := range buckets {
		if bucket.Name == path {
			if _, ok := bucketObjectMap[bucket]; ok {
				return &Data{IsBucket: false, objects: bucketObjectMap[bucket]}
			}
		}
	}

	return nil
}
