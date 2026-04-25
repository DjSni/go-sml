package sml

type GetProfileListRequest = GetProfilePackRequest

func GetProfileListRequestParse(buf *Buffer) (GetProfileListRequest, error) {
	return GetProfilePackRequestParse(buf)
}
