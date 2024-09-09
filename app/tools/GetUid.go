package tools

import "github.com/bwmarrin/snowflake"

var snowNode *snowflake.Node

func GetUid() int64 {
	if snowNode == nil {
		snowNode, _ = snowflake.NewNode(1)
	}
	return snowNode.Generate().Int64()
}
