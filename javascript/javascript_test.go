package javascript

import (
	"github.com/team-ide/go-interpreter/node"
	"testing"
)

const javaScriptCode = `

var a = 1;
var b = 1;
a=1
b=2
aa=[1,2,3,""]
s+=s
s++
s--
b=j++;
b=j--;
s=s*s-1+4*s
let a = ''
let a = null
a=a;
;;
a=[{a:1},[{s:1,b:2}]]
var aFunc = function(){

}

aaa.each(()=>{

})

for(;;){
}

for(a=1;a<10;a++){
}

for(a in aa){
}

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
var x = 5;
const PI = 3.14;

var num = 5;
var str = "Hello";
var bool = true;
var arr = [1, 2, 3];
var obj = {name: "John", age: 30};
function sayHello() {
  console.log("Hello");
}

var z = x + y;
var a = x > y;
var b = x || y;

if (x > 5) {
  console.log("x is greater than 5");
} else {
  console.log("x is less than or equal to 5");
}

for (var i = 0; i < 10; i++) {
	console.log(i);
	if(i==10){
		break;
	}else{
		continue;
	}
}

while (x < 5) {
  console.log(x);
  x++;
}


function add(x, y) {
  return x + y;
}

var subtract = function(x, y) {
  return x - y;
};

var obj = {name: "John", age: 30};
obj.name = "Mary";
console.log(obj.age);

class Person {
  constructor(name, age) {
    this.name = name;
    this.age = age;
  }

  sayHello() {
    console.log("Hello, my name is " + this.name);
  }
}

var person = new Person("John", 30);
person.sayHello();

setTimeout(function() {
  console.log("Hello, world!");
}, 1000);

fetch("https://api.example.com/data")
  .then(function(response) {
    return response.json();
  })
  .then(function(data) {
    console.log(data);
  })
  .catch(function(error) {
    console.log(error);
  });

async function getData() {
  try {
    const response = await fetch("https://api.example.com/data");
    const data = await response.json();
    console.log(data);
  } catch (error) {
    console.log(error);
  }
}

// Shape - superclass
function Shape() {
    this.x = 0;
    this.y = 0;
}

// superclass method
Shape.prototype.move = function(x, y) {
    this.x += x;
    this.y += y;
};

// Rectangle - subclass
function Rectangle() {
    Shape.call(this);
}

// subclass extends superclass
Rectangle.prototype = Object.create(Shape.prototype);
Rectangle.prototype.constructor = Rectangle;

var rect = new Rectangle();

// Shape - superclass
function Shape() {
    this.x = 0;
    this.y = 0;
	super()
	super.x =1
	super.x.s =1
	super[0]
}

// superclass method
Shape.prototype.move = function(x, y) {
    this.x += x;
    this.y += y;
};

// Rectangle - subclass
function Rectangle() {
    Shape.call(this);
}

// subclass extends superclass
Rectangle.prototype = new Shape(ss);

var rect = new Rectangle( );
`

func TestJavaScript(t *testing.T) {
	tree, err := Parse(javaScriptCode)
	if tree != nil {
		node.OutTree(javaScriptCode, tree)
	}
	if err != nil {
		panic(err)
	}
}
