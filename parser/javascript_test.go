package parser

import (
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-interpreter/language"
	"github.com/team-ide/go-interpreter/node"
	"testing"
)

const javaScriptCode = `

var a = 1;
var b = 1;
var aFunc = function(){

}
a=1
b=2
aa=[1,2,3,""]

aaa.each(()=>{

})

for(;;){
}

for(a=1;a<10;a++){
}

for(a in aa){
}

s+=s
s++
s--
s=s*s-1+4*s
let a = ''
let a = null

a=a;
;;
a=[{a:1},[{s:1,b:2}]]
try{
}catch(e){
}finally{
}

function a(a,b){

return aaa;
}
for (element of iterable) {
    // body of for...of
}
debugger;

let aa = async (a,b)=>{
await aa(a,c);
}

function Cat(){
    this.name = "大毛";
}
var Cat = {
    name:'大毛',
    makeSound:function(){
        alert('喵喵喵');
    }
}

var Cat = {
    createNew:function(){
       var cat = {};
       cat.name = '大毛';
       cat.makeSound = function(){
           alert('喵喵喵');
       }
       return cat;
    }
}
if (a==1 && a==2){
	a = 1
}else if(a==2){
	a = 1

}else{
	a = 1
	throw new Error('a')
}
switch(aa){
case "1":
		break;
case "2":
		break;
default:
		break;
}
while (i < 10) {
    text += "数字是 " + i;
    i++;
}
do {
    text += "The number is " + i;
    i++;
 }
while (i < 10);

let aa = /xx^/
aa.test('')
const PI = 3.141592653589793;
PI = 3.14;      // 会出错
PI = PI + 10;   // 也会出错
var person = {
  firstName: "Bill",
  lastName : "Gates",
  id       : 678,
  fullName : function() {
    return this.firstName + " " + this.lastName;
  }
};
with(location){
    let qs = search.substring(1);
    let hostName = hostname;
    let url = href;
    console.log('with:',qs,hostName,url);
}

class Rect {
    // 固定名称--称为构造方法
    // 当 new 类名 时, 就会触发此方法
    constructor(width, height) {
        console.log('constructor 被触发!');
        this.width = width
        this.height = height
    }
    // 省略 prototype, function ...
    // 会自动加入到原型中
    area(){
        return this.width * this.height
    }
 
    zc(){
        return (this.width + this.height) * 2
    }
}
 
var r1 = new Rect(10, 20)
console.log(r1);
 
console.log(r1.area());
console.log(r1.zc());
`

func TestJavaScript(t *testing.T) {
	tree, err := Parse(javaScriptCode, &language.JavaScriptSyntax{})
	if err != nil {
		panic("parser.Parse error:" + err.Error())
	}
	outTree(javaScriptCode, tree)

}

func outTree(code string, tree *node.Tree) {
	for _, one := range tree.Children {
		bs, _ := json.Marshal(one)
		fmt.Println("tree one start:", one.Start()-1, ",end:", one.End()-1, ",info:", string(bs))
		fmt.Println(code[one.Start()-1 : one.End()-1])
	}
}
