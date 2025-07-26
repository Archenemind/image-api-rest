package images

type Image struct {
	Id, Username, ImageName, Size, Format string
}

var Images = []Image{
	{Id: "1", ImageName: "image1", Size: "1", Format: "png"},
	{Id: "2", ImageName: "image2", Size: "2", Format: "png"},
	{Id: "3", ImageName: "image3", Size: "3", Format: "png"},
	{Id: "4", ImageName: "image4", Size: "4", Format: "png"},
	{Id: "5", ImageName: "image5", Size: "5", Format: "png"},
}
