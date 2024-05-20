package controller

import (
	"backend/pkg"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 保存了农产品上链与查询的函数

func Uplink(c *gin.Context) {
	userType, _ := c.Get("userType")
	file, _ := c.FormFile("file")
	// 从表单中获取文件
	// 与userID不一样，取ID从第二位作为溯源码，即18位，雪花ID的结构如下（转化为10进制共19位）：
	// +--------------------------------------------------------------------------+
	// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
	// +--------------------------------------------------------------------------+
	farmer_traceability_code := pkg.GenerateID()[1:]
	args := buildArgs(c, farmer_traceability_code)
	traceability_code := args[1]
	res, err := pkg.ChaincodeInvoke("Uplink", args)

	if err != nil {
		c.JSON(200, gin.H{
			"message": "uplink failed" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":              200,
		"message":           "uplink success",
		"txid":              res,
		"traceability_code": args[1],
	})
	if userType == "种植户" {
		// 创建保存文件的路径
		dst := "dist/upload/" + traceability_code + ".jpg"
		// 保存文件
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(200, gin.H{
				"message": "Unable to save the file",
			})
			return
		}
	}
}

// 获取图片
// func GetPic(c *gin.Context) {
// 	tag := c.PostForm("tag")
// 	filePath := "http://127.0.0.1/9090:images/" + tag + ".jpg" // 修改为实际存储图片的路径格式

// 	// 检查文件是否存在
// 	_, err := os.Stat(filePath)
// 	if os.IsNotExist(err) {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"code":     200,
// 		"imageUrl": filePath,
// 	})
// }

// 获取农产品的上链信息
func GetFruitInfo(c *gin.Context) {
	traceability_code := c.PostForm("traceability_code")
	res, err := pkg.ChaincodeQuery("GetFruitInfo", traceability_code)
	path := "./dist/upload/" + traceability_code + ".jpg"
	if err != nil {
		c.JSON(200, gin.H{
			"message": "查询失败：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
		"path":    path,
	})

}

// 获取用户的农产品ID列表
func GetFruitList(c *gin.Context) {
	userID, _ := c.Get("userID")
	res, err := pkg.ChaincodeQuery("GetFruitList", userID.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

// 获取所有的农产品信息
func GetAllFruitInfo(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetAllFruitInfo", "")
	fmt.Print("res", res)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

// 获取农产品上链历史
// func (s *SmartContract) GetFruitHistory(ctx contractapi.TransactionContextInterface, traceability_code string) ([]HistoryQueryResult, error) {
func GetFruitHistory(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetFruitHistory", c.PostForm("traceability_code"))
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed" + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

func buildArgs(c *gin.Context, farmer_traceability_code string) []string {
	var args []string
	userID, _ := c.Get("userID")
	userType, _ := pkg.ChaincodeQuery("GetUserType", userID.(string))
	args = append(args, userID.(string))
	fmt.Print(userID)
	// 种植户不需要输入溯源码，其他用户需要，通过雪花算法生成ID
	if userType == "种植户" {
		args = append(args, farmer_traceability_code)
	} else {
		args = append(args, c.PostForm("traceability_code"))
	}
	args = append(args, c.PostForm("arg1"))
	args = append(args, c.PostForm("arg2"))
	args = append(args, c.PostForm("arg3"))
	args = append(args, c.PostForm("arg4"))
	args = append(args, c.PostForm("arg5"))
	return args
}
