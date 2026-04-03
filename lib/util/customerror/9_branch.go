package customerror

const featBranchNum = 9

var (
	BRANCH_DUPLICATE_NAME                        = New(featBranchNum, 0, "Branch name already exists")                                 // 09000
	BRANCH_GEOM_POINT_IS_NOT_RELATED_TO_LOCATION = New(featBranchNum, 1, "Branch geom point is not related to the specified location") // 09001
)
