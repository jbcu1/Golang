package recon

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func OpenAndReadTrancoCSV() []string{
	fileName:=GetArchive()
	file, err:=os.Open(fileName)
	if err!=nil{
		fmt.Errorf("Sometring went wrong %v\n", err)
	}
	defer file.Close()
	reader:=csv.NewReader(file)
	reader.FieldsPerRecord =2
	linksArray:=make([]string,0)
	for{
		record, err:=reader.Read()
		if err!=nil {
			break
		}
		linksArray = append(linksArray, record[1])
	}

	return linksArray
}


func GetArchive() string{
	fmt.Println("Start copy")
	fileName,folderPath,err:=downLoadArchive("tl.zip","http://tranco-list.eu/download_daily/P5PJ")
	if err!=nil{
		fmt.Errorf("Something went wrong %s.\n",err)
	}
	files, err:=unzipArchive(fileName,folderPath)
	if err!=nil{
		fmt.Errorf("Something went wrong %s.\n",err)
	}
	removeFile(fileName)

	return files[0]

}


var trancoFolder = "tranco_result"

//Delete file
func removeFile(pathToFile string)error{
	err:=os.Remove(pathToFile)
	if err!=nil{
		return err
	}
	return nil
}

func downLoadArchive(fileName string, url string) (string,string,error){
	currPath,_:=os.Getwd()

	resp, err:=http.Get(url)
	if err!=nil{
		fmt.Errorf("Someting went wrong, when get request. %s\n", err)
	}
	defer resp.Body.Close()


	os.MkdirAll(currPath+"/pars_result/"+trancoFolder,0777)
	out, err:=os.Create(currPath+"/pars_result/"+trancoFolder+"/"+fileName)
	if err!=nil{
		fmt.Errorf("Someting went wrong, when create file. %s\n", err)
	}

	zipParh:=currPath+"/pars_result/"+trancoFolder+"/"+fileName
	folderPath:=currPath+"/pars_result/"+trancoFolder
	defer out.Close()

	_,err=io.Copy(out,resp.Body)



	return zipParh,folderPath,nil
}


func unzipArchive(src string, dest string) ([]string,error){
	var filenames []string
	r,err:=zip.OpenReader(src)
	if err!=nil{
		return filenames,err
	}
	defer r.Close()

	for _,f:=range r.File{
		fpath:=filepath.Join(dest,f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)){
			return filenames,fmt.Errorf("%s: illegal file path", fpath)
		}
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir(){
			os.MkdirAll(fpath,0777)
		}
		if err = os.MkdirAll(filepath.Dir(fpath),0777); err!=nil{
			return filenames,err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err:=f.Open()
		if err!=nil{
			return filenames,err
		}

		_, err=io.Copy(outFile,rc)

		outFile.Close()
		rc.Close()

		if err!=nil{
			return filenames,err
		}
	}


	return filenames,nil
}