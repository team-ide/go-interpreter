package javascript

/**
符号

算术运算符：+、-、*、/、% 分别表示加、减、乘、除、取模运算。

赋值运算符：=、+=、-=、*=、/=、%= 分别表示赋值、加等于、减等于、乘等于、除等于、取模等于运算。

比较运算符：==、===、!=、!==、>、>=、<、<= 分别表示等于、恒等于、不等于、不恒等于、大于、大于等于、小于、小于等于运算。

逻辑运算符：&&、||、! 分别表示与、或、非运算。

位运算符：&、|、^、~、<<、>>、>>> 分别表示按位与、按位或、按位异或、按位取反、左移、右移、无符号右移运算。

三目运算符：? : 表示条件运算符，用于简单的条件判断。

其他运算符：typeof、delete、in、instanceof、new、void、yield 等，用于类型检查、删除对象属性、判断属性是否存在、判断对象是否为某个类型、创建对象实例、计算表达式并返回 undefined、生成迭代器等。

JavaScript 中的符号和运算符都有其特定的用法和优先级，开发者需要熟练掌握才能正确地使用它们。
*/

/**
运算符

算术运算符：+，-，*，/，%，++，--

关系运算符：==，===，!=，!==，>，<，>=，<=

逻辑运算符：&&，||，!

位运算符：&，|，^，~，<<，>>，>>>

赋值运算符：=，+=，-=，*=，/=，%=，<<=，>>=，&=，^=，|=

三元运算符：? :

instanceof 运算符：用于判断一个对象是否是某个类的实例。

in 运算符：用于判断一个对象是否包含指定的属性。

delete 运算符：用于删除对象的属性或数组中的元素。

typeof 运算符：用于返回一个值的类型。

void 运算符：用于指定表达式没有返回值。

除此之外，JavaScript 还有一些特殊的运算符，如：

箭头函数运算符：=>，用于描述箭头函数的参数和方法体。

条件运算符：用于描述一组表达式与值的关系，如 switch 语句中的 case 关键字。

点运算符和中括号运算符：用于访问对象的属性或方法。

这些运算符可以用于不同的数据类型，如整数、浮点数、布尔值、字符、字符串和对象等。
*/

/**
关键字

保留字：break、case、catch、class、const、continue、debugger、default、delete、do、else、export、extends、finally、for、function、if、import、in、instanceof、new、return、super、switch、this、throw、try、typeof、var、void、while、with。

严格模式下的保留字：implements、interface、let、package、private、protected、public、static、yield。

未来的保留字：await、enum。

这些关键字都有特定的用途和语法规则，开发者在使用时需要注意避免与关键字冲突。
*/

/**
基础类型

Undefined：表示未定义或未初始化的值。

Null：表示空对象指针。

Boolean：表示布尔值，即 true 或 false。

Number：表示数值，包括整数和浮点数，还包括特殊的值，如 Infinity、-Infinity 和 NaN。

String：表示字符串，由一组 16 位 Unicode 字符序列组成。

Symbol：表示唯一的标识符。

BigInt：表示大整数，用于处理超出 JavaScript Number 类型范围的整数。

需要注意的是，JavaScript 中的基础类型（除了对象类型）都是不可变的，也就是说，一旦创建就不能再修改其值，而是会创建一个新的值。
*/
