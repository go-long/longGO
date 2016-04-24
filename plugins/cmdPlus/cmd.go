package cmdPlus

import (
	"os/exec"
	"bufio"
"io"
	"regexp"
	"github.com/hoysoft/JexGO/logger"
	"fmt"
)

type CmdPlus struct {
	Cmd        *exec.Cmd
	OutPutCallback OutPutCallbackFunc
	regexpKeys []string
	TriggerKeyCallback TriggerKeyCallbackFunc //关键字触发回调
}

type FinishCallbackFunc func (error);
type TriggerKeyCallbackFunc func (map[string]string);
type OutPutCallbackFunc func (string);

func NewCmdPlus(cmd *exec.Cmd )*CmdPlus{
	c:=&CmdPlus{Cmd:cmd}
	return c
}

func (this *CmdPlus)SetTriggerRegexpKeys(m ...string)*CmdPlus {
	this.regexpKeys=m
	return this
}

func (c *CmdPlus)Exec() error {

	//显示运行的命令
	logger.Info("cmd.Args:",c.Cmd.Args)
	logger.Info("cmd WorkDirectory:",c.Cmd.Dir)
	stdout, err := c.Cmd.StdoutPipe()

	if err != nil {
		return   fmt.Errorf("RunCommand: cmd.StdoutPipe(): %v", err)

	}

	stderr, err := c.Cmd.StderrPipe()

	if err != nil {
		return     fmt.Errorf("RunCommand: cmd.StderrPipe(): %v", err)
	}
	c.Cmd.Start()

	go func(){
		c.readerLine(bufio.NewReader(stdout))
	}()
	go func(){
		c.readerLine(bufio.NewReader(stderr))
	}()


	err2:=c.Cmd.Wait()
	//if err2!=nil{
	//	fmt.Println("return err:",err2,c.Cmd.ProcessState)
	//}
	return err2
}

func (c *CmdPlus)readerLine(reader *bufio.Reader){
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		//正则表达式提取关键字触发回调
		c.regexpTriggerKeys(line)
		if c.OutPutCallback!=nil{
			c.OutPutCallback(line)
		}
	}
}

//正则表达式提取关键字触发回调
func (this *CmdPlus)regexpTriggerKeys(line string){
	if this.TriggerKeyCallback==nil {return }
	for _,v:=range this.regexpKeys{
		var digitsRegexp = regexp.MustCompile(v)
		m:=FindStringSubmatchMap(digitsRegexp,line)
		if m!=nil && len(m)>0  {
			this.TriggerKeyCallback(m)
		}
	}

}

func FindStringSubmatchMap(r *regexp.Regexp,s string) map[string]string{
	captures:=make(map[string]string)

	match:=r.FindStringSubmatch(s)
	if match==nil{
		return captures
	}

	for i,name:=range r.SubexpNames(){
		//Ignore the whole regexp match and unnamed groups
		if i==0||name==""{
			continue
		}

		captures[name]=match[i]

	}
	return captures
}