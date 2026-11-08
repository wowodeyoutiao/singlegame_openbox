package excel

import (
	"fmt"
	"nanoserver/global"

	"os"
	"path"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func LoadExcel(fileName string) ([][]string, error) {
	// 获取上传文件

	// 设置文件需要保存的指定位置并设置保存的文件名字
	dst := path.Join(global.BasePath+"/configs/gameconfig/", fileName)
	//dst := path.Join("./configs/gameconfig/", fileName)
	xlsx, err := excelize.OpenFile(dst)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 获取excel中具体的列的值
	rows := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))

	return rows, err
}
