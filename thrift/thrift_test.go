package thrift

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"testing"
)

const thriftCode = `
include "a.thrift"

namespace java com.thrift.service
namespace go com.thrift.service
namespace cpp com.thrift.service

/**
 * 异常信息
 */
exception XException{
	/**
     * 错误码
	 */
    1: optional i32 code;
	/**
     * 错误信息
	 */
    2: optional string msg;
}


/**
 * 对象
 */
struct Obj {
  1: i32 num1 = 0;		//默认值
  2: i32 num2;
  3: Operation op;		//可以嵌套其他类型
  4: optional string comment;	//可选字段
  5: list<i32> l		//list
  6: map<i32,string> m	//map
  7: set<string> s	//set
  8: bool vBool	//bool
  9: byte vByte	//byte
  10: i16 vI16	//i16
  11: i64 vI64	//i64
  11: double vDouble	//double
  12: string vString	//string
  13: Obj2 vObj2	//Obj2
}

struct Obj2 {
  1: i32 num1 = 0;		//默认值
}

/**
 * 枚举
 */
enum Operation {
  ADD = 1,
  SUBTRACT = 2,
  MULTIPLY = 3,
  DIVIDE = 4
}

/**
 * 响应
 */
struct XResponse{
	/**
     * 状态码
	 */
    1: optional i32 code;
	/**
     * 状态信息
	 */
    2: optional string msg;
}


/**
 * 服务
 */
service Service {
	/**
	 * 方法
	 * 参数：
	 *	    Obj request
	 * 返回：
	 *		XResponse res
	 */
	XResponse method1(1: Obj request);

}


include "aphead.thrift"
namespace java com.vrv.im.service
namespace cpp com.vrv.im.service

struct LoginBean {
  /*账号、邮箱、手机号、QQ号、（湖北公安code）*/
  1: string loginAccount; 
  /*1：PC；2：手机；3：pad；4：网页*/
  2: byte deviceType;   
  /*接入点IP*/
  3: string ip;
  /*接入点端口*/
  4: i32 port;
  /*用户类型：1：手机号 2：第三方关联,3：邮箱,4：用户ID,99: SAML认证用户,6：豆豆号,7:使用自定义拓展字段来进行登陆, 8: 第三方授权登录 -1：传用户ID但票验证 100湖北公安auth登录*/
  5: byte userType;
  /*mac地址*/
  6: string macCode;
  /*设备信息*/
  7: string deviceInfo;
  /*-1：离线，1:在线；2:隐身,3:忙碌,4:离开,5:请勿打扰,6:想聊天*/
  8: byte flag;
  /* 客户端IP */
  9: optional string clientIP;
  /* 是否已验证，true已验证(针对resSign=8|9的时候用) */
  10: optional bool verified;
  11: i64 SDKID;
  /*auth回掉地址*/
  12:optional string redirectUri;
  /*设备类型 品牌主要是推送用 0苹果 1 华为 2 google 3 中电科*/
  13:optional i32 deviceType2;
  /*密码(华数传媒)*/
  14:optional string pwd;
  /*客户端版本号*/
  15:optional string clientVersion;
  /*用于传递三方认证所需参数（JSON格式）*/
  16:optional string thirdAuthParamJson;
  /* 苹果ID */
  17: optional string appleID;
  /* 推送token */
  18: optional string token;
  /* IOS推送token */
  19: optional string voipToken;
}
struct TicketBean {
  /*1：成功；2：账号、密码不匹配；3：已经存在登录；4：无此ticket, 5:账户被锁定，6:账户被冻结，7:账户已停用， 8：非常用登录设备, 9：非常用登录地址, 10:需要修改密码, 11:需要登录验证 ,19:设备被锁定*/
  1: byte resSign; 
  /*成功时返回新建ticket；不匹配时返回空；已经存在登录的时候返回存在的ticket。*/
  2: optional string ticket;
  /*成功时返回connectID；不匹配时返回空；已经存在登录的时候返回连接connectID*/
  3: optional i64 connectID;
  /*成功时返回userID；不匹配时返回空；已经存在登录的时候返回userID*/
  4: optional i64 userID;
   /*成功时返回area；不匹配时返回空；已经存在登录的时候返回area*/
  5: optional string area;
  /* 账户被锁定时，返回剩余锁定时长 ，毫秒*/
  6: optional i64 remainLockTime;
  7:i64 SDKID;
  /* 登录设备*/
  8:optional string device;
  /* 登录位置*/
  9:optional string location;
  /* 用户注册时间*/
  10:optional i64 regTime;
  /* 用于传递三方认证返回信息*/
  11:optional string thirdAuthResultJson;
  /**
    * 扩展字段
    *{
    *	  "customerService":["key1","key2"], //我的客服 key1为parentId=20001的拓展字典的key
    *	  "crtInfo": {
    *					 "url":"vrv.linkdood.cn", //证书续费地址，只有绑定管理员的账户才会返回该字段
    *					 "serverTime":1213454435000, //服务期当前时间（时间戳，毫秒）
    *					 "crtBTime":1213454435000, //证书有效期开始时间（时间戳，毫秒）
    *					 "crtETime":1213454435000 //证书有效期到期时间（时间戳，毫秒）
    *				  },
    *	  "MHInfo":{
    *		   "adminUrl":"/mbox" //魔盒管理中心
    *	  },
    *	  "verifiedInfo":{
    *		   "region":"CN", // COM海外 , 其他中国
    *		   "personSize":10 // int 手机号注册用户数
    *	  },
    *	  "commInfo":{
    *			"authority":1 // int,是否有实名担保标识字段; 0 没有; 1 有
    *	  },
    *     "userType":0 // int,用户类型,字段不存在或为0表示正式用户,为9表示试用用户
    * }
    */
   12: optional string exinfo;

   /** 剩余登录次数*/
   13: i32 remainingLoginTimes;
}

/**
* 登录状况
*/
struct LoginState{
    /* 用户ID */
    1: i64 userID;
    /* 状态: 1正常, 2异常 */
    2: byte status;
    /* 更新时间 */
    3: i64 updateTime;
    /* 异常登录记录, 只有在异常状态的时候返回 */
    4: optional list<aphead.LoginRecord> records;
}

/**
 * 常用登录设备
 */
struct LoginDevice{
	/* 用户ID */
	1: i64 userID;
	/* 设备类型 */
	2: byte deviceType;
	/* 设备信息 */
	3: string deviceInfo;
	/* mac地址 */
	4: string macCode;
	5: i64 SDKID;
}

/**
 * 票据参数
 */
struct LoginTicket{

	/* 票据ID */
	1: string vrvID;
	/* 票据信息 */
	2: string vrvTicketInfo;
	
}

/**
 * 设备锁定bean
 */
struct DeviceInfoBean{   
	/* 用户们的ID */
	1: string userIds;
	/* 设备类型  1：PC；2：android；3：ios */
	2: byte deviceType;
	/* 设备信息 */
	3: string deviceInfo;
	/* mac地址 */
	4: string macCode;
	/* 设备锁定时间 */
	5: i64 lockDeadline;
	/* 用户们的名称 */
	6: string userName;
}

/**
 * 新设备锁定bean
 */
struct DeviceLock{
	/* 用户们的ID */
	1: string macCode;
	/* 设备信息 */
    2: string deviceInfo;
    /* 状态 1锁定 2解锁 */
    3: i32 state;
	/* 失效时间 */
	4: i64 lockDeadline;
	/* 创建时间 */
    5: i64 createTime;
    /* 失效时间 */
    6: i64 updateTime;

}


/**
 * 设备锁定分页查询
 */
struct DeviceInfoPage{
    /*页码*/
	1:optional i32 pageNum;
	/*页长*/
	2:optional i32 pageSize;
	/*总数*/
	3:optional i32 count;
	/*总页数*/
	4:optional i32 totalPage;
	/*操作日志列表*/
	5:optional list<DeviceLock> ams;
}

/**
 * 用于接收adapter-login包的doLoginPre方法的返回结果AdapterLoginPreResult
 */
struct ServerLoginPreResult{
	/*返回码*/
	1: i32 resultCode;
	/*返回信息*/	
	2: string resultMsg;
	/*返回数据*/
	3: string resultData;
}

/**
 * 用于getUserOnlineByLogin的结果
 */
struct UserOnlineResult{
  /*返回码*/
  1: i32 resultCode;
  /*返回数据*/
  2: list<aphead.OnLineUserBean> onlineBeans;
  /*用户id*/
  3: i64 userId;
}

service LoginService {
  /**
   ******************************登录请求创建票/验票服务****************************
   *用户类型非-1的时候逻辑
   *客户端调用的时候，发送登录信息，如果是有登录，提示是否强制登录，如果选择是
   *就发送强制登录标识isForceKick，调用强制登录，另外登录的状态，登录成功后再发
   *修改状态信息（onlineService的updateOnlineState）。在AP接入点返回成功信息  
   *后，正常登录的时候需要将userID通知路由服务(notifyChange)，连接改变。      
   *由路由通知chatservice。                                                  
   *            客户端通过接入点请求                                          
   *                   |                                                      
   *                   |                                                      
   *     先判断用户密码是否匹配(用户服务获取密码)，不匹配直接返回，           
   *             匹配得到userID,往下走                                        
   *                   |                                                      
   *                   |                                                      
   *         判断相同userId，相同的设备，有没有ticket                         
   *             |                                          |                 
   *         没有 |                                         有|                              
   *     生成ticket，并且                           判断是不是isForceKick     
   *   去onlineService的Registry注册生成connectId   如果是2,直接返回已经  
   *	返回ticket和connectId			              存在登录标识。             
   *     去onlineService注册生成connectId           如果是1，如果是就调用  
   *				                              onlineService的kickRegistry 
   *						                      删除原来用户，相同设备下的 
   *						                      connect信息，注册生成新的  
   *						                      connectId，并返回connectId 
   *--------------------------------------------------------------------------
   *用户类型为-1逻辑，票验证逻辑
   *客户端通过验票登录，cookie处理
   *       loginBean 传递loginAccount 传递userID   
   *       isForceKick不做处理
   *       获取该设备类型的ticket，如果不存在，返回无效票（客户端出登录界面）
   *                    |
   *                    |
   *       判断该票是否一致，不一致，mac一致返回无效票（客户端出登录界面）
   *       mac地址不一致，返回有其他设备登录
   *                    |
   *       获取新的连接ID，返回验证成功
   *
   */
  TicketBean createTicket(1: LoginBean lb, 2: string pwd,3: byte isForceKick)
   /**
       ******************************校验用户/返回用户在线信息****************************
       *          根据账号类型及密码校验用户是否存在
       *          不存在返回错误码,登录失败
       *          存在则通过用户ID查询该用户所有在线信息
       *          异常返回服务器异常
       *
       */
      UserOnlineResult getUserOnlineByLogin(1: LoginBean lb, 2: string pwd)
  /**
   ******************************删除票服务************************************
   * 客户端正常退出（客户端异常中断连接，ap直接去在线服务调用删除在线信息）                                    
   *            客户端主动退出通过接入点请求                       
   *                   |                                                     
   *                   |                                                    
   *     调用onlineService删除连接ID服务                                      
   *                   |                                                     
   *                   |                                                     
   *          删除本ticket信息                                                
   */
   void deleteTicket(1:string ticket, 2:i64 userId, 3:byte deviceType, 4:i64 SDKID)
   
  /**
   ******************************客户端删除pc票服务************************************
   *  客户端通过ticketId删除pc票据,并使pc客户端下线                                                                                                                       
   */
   aphead.Result deleteTicketWithoutTicket(1:byte deviceType,2:i64 userId,3:i64 SDKID)
      
   
  /**
   ******************************删除用户所有票服务************************************
   * 用户账号冻结时调用                                                                                      
   *      删除该用户所有ticket信息                                             
   */
   void clearUserAllTicket(1:i64 userId,2:i64 SDKID)  
  
  /**
   ******************************网页登录****************************
   *网页登录，默认强踢，目前是单点登录用。
   *返回登录结果 ip和port不传递
   */  
   
   
  TicketBean webLogin(1: LoginBean lb, 2: string pwd)
  
  /**
   ******************************SSO登录****************************
   * SSO登陆接口，
   * 返回登录结果：票据ticket和用户userID
   */  
  TicketBean ssoLogin(1: LoginBean lb, 2: string pwd)
  
  /**
   ******************************SSO登录校验****************************
   * SSO登陆验证接口
   * 返回执行结果
   */  
  aphead.Result ssoLoginTicketVerify(1: LoginTicket loginTicket)
    
  /**
   ******************************SSO票据删除****************************
   * SSO票据删除
   * 返回执行结果
   */  
  aphead.Result ssoLoginTicketRemove(1: LoginTicket loginTicket)
  
  /**
   ******************************获取网页访问临时key****************************
   *返回创建临时key，并放入缓存中，临时key，每个用户就产生一个，新的顶掉老的
   *返回临时key
   */  
  string getClientKey(1: i64 userID, 2: string ticket,3:i64 SDKID)
  
  /**
   ******************************linkdood内部跳转应用获取授权码****************************
   *返回创建临时Code，调用platform接口存储Code和UserID关系
   *返回临时Code
   */  
  string getLoginAuthCode(1: aphead.LoginAuthParam loginAuthParam)
  
  /**
   *******************************验证key***********************************
   *解析clientKey看是否存在，是否过期，
   *清除该用户的clientKey，一个key只用一次。
   *返回userID
   */  
  i64 verifyClientKey(1: string clientKey,2:i64 SDKID)
  /**
   *******************************上报通电状态***********************************
   *status 1:充电，2:非充电
   *
   */  
   void reportPlugInStatus(1:i64 userID,2:byte deviceType,3:string macCode,4:byte status,5:i64 SDKID)
  /**
   *******************************查询通电状态***********************************
   *返回 1:充电，2:非充电
   *
   */ 
  byte getPlugInStatus(1:i64 userID,2:byte deviceType,3:string macCode,4:i64 SDKID)
  
  /**
   *******************************获取最近登录历史***********************************
   * 场景：WEB调用
   * 参数：userID 用户ID
   * 返回：最近50次登录历史记录
   */
  list<aphead.LoginRecord> getLoginRecentHistory(1:i64 userID,2:i64 SDKID)	
  
  /**
   *******************************获取常用登录设备列表***********************************
   * 场景：WEB调用
   * 参数：userID 用户ID
   * 返回：常用登录设备列表
   */
  list<LoginDevice> getLoginDevices(1:i64 userID,2:i64 SDKID)

  /**
   ********************获取登录情况****************************
   * 场景: WEB调用
   * 参数: userID 用户ID
   * 返回: 登录异常情况, 异常状态时同时返回最近5条异常登录记录
   */
  LoginState getLoginState(1:i64 userID, 2:i64 SDKID)

  /**
   * 验证用户票是否存在
   * 场景: 服务|公众平台
   * 参数: userID 用户ID
   *      ticket 用户票
   * 返回: true:存在, false:不存在
   */
  bool validTicket(1: i64 userID, 2: string ticket,3:i64 SDKID)
  /**
   * 统计登陆人员总数
   * 入参data:表示要查询的日期，精确到天，例如：2016-11-21
   *                             目前data传递  表示按日期查询当天的活跃用户数量
   *					     	   目前data不传   表示登录用户总数量
   *
   * 返回值：表示用户数数量
   */
  i64 queryLoginUserTotalCount(1:string data)
  
   /**
   *******************************查询用户最后一次登录信息***********************************，
   *入参：用户ID
   *返回userID
   */  
  i64 getLoginLangType(1: i64 userID,2:i64 SDKID)
  
  /**
   * 根据客户端类型获取客户端各版本的使用量
   * 入参：
   * 	clientType（1.PC 2.android 3.iOS 4.元心 5.深度 6.麒麟 7.其他客户端）
   * 返回：
   * 	{2.0.36=3, other=544}
   */
  map<string,i64> getAllClientVersionUserCounts(1:byte clientType)

    /**
      * 获取authCode认证
      * resultCode:
      * 100008000 服务异常
      * 100008002 获取临时授权码失败
      * 100008001 应用不合法
      * 0 成功
      **/
     aphead.AuthCodeResult getAuthCode(1: aphead.AuthCodeParam authCodeParam);

     /**
      * 授权认证二维码校验
      * resultCode:
      * 100008000 ->服务异常
      * 100008001 ->公众号不合法
      * 100008003 ->二维码失效
      * 0         ->成功
      **/
     aphead.QrcodeValidateResult qrcodeValidate(1: aphead.QrcodeValidateParam validateParam);

     /**
      * 允许授权登陆
      * 返回值：
      * 100008000 -> 服务异常
      * 100008003 -> 二维码失效
      * 0         -> 成功
      **/
     i32 qrcodeLogin (1:string qrcode);
	 
	 /**
   * 根据用户ID解锁相关的所有设备
   * 入参：
   * 	
   * 返回：
   * 	
   */
  void unlockDeviceByUserID(1: i64 userID)

  	/**
     * 获取用户登录的设备信息
     * 
     * 返回
     **/
    list<byte> getDeviceType(1: i64 userID,2: i64 SDKID,3: byte userType);

    /**
     * 客户端扫码成功接口
     * @param qId 二维码唯一标识ID
     * @param userId 用户唯一标识ID
     * @param serverUrl 服务器地址
     *
     * 返回
     **/
    aphead.QrCodeResultBean scanSuccess (1:i64 qId,2:i64 userId,3:string serverUrl)

     /**
     * 客户端扫码登录接口
     * @param qId 二维码唯一标识ID
     * @param userId 用户唯一标识ID
     * @param serverUrl 服务器地址
     *
     * 返回
     **/
    aphead.QrCodeResultBean qCodeLogin (1:i64 qId,2:i64 userId,3:string serverUrl);
    
    /**
	 * 透传接口：用于第三方认证前置的校验业务
	 * @param method 请求方法
	 * @param paramJson 请求参数
	 */
	ServerLoginPreResult doLoginPre(1:string method,2:string paramJson);

   /**
     * 用户信息字段验证接口
     * @param accountType 用户账号信息
     * @param userFieldJson 用户验证字段信息(JSON格式)，示例{"base_name": "value1","account_10": "value2","extend_49114af6d8d64aa7b9ff7d79ab7be137": "value3"}
     *
     * 返回 100008004用户不存在 100008005验证用户信息失败 100008006验证用户信息成功
     */
     aphead.Result verifyUserInfoField(1:aphead.AccountType accountType,2:string userFieldJson,3:i64 SDKID);
     
     /**
      * 获取授权设备
      */
     list<aphead.LoginRecord> getLoginDeviceRecords(1:i64 userID, 2:i64 SDKID);
     
     /**
      * 移除授权设备
      */
     aphead.Result deleteLoginDeviceRecords(1:i64 userID, 2:list<string> macCodeList, 3:i64 SDKID);

	  /****************************************设备信息分页查询*********************************
	  * 分页查询菜单
	  *  limit=每页的条数
	  *  page=当前页码
	  *  map 参数 1.name 精确查询
	  * 
	  */
	DeviceInfoPage querylockDevicePage(1:map<string,string> wheremap,2:i32 page,3:i32 limit);
	
	/**
      * 解锁设备
      */
     aphead.Result unlockDevice(1:string macCode,2:i64 SDKID);

     /**
       * 删除锁定设备记录信息
       */
      aphead.Result deleteLockDeadline(1:string macCode,2:i64 SDKID);


      /****************************************PC扫码登录 start*********************************/

    /**
     * 获取PC登录二维码信息
     * @param elogo 服务器elogo
     * @param loginData PC端登录信息
     * @param SDKID
     *
     */
     aphead.PCLoginQrCodeResult getPCLoginQrCode(1:string elogo, 2:string loginData,3:i64 SDKID);

     /**
      * 移动端扫PC端登录二维码接口
      * @param elogo 服务器elogo
      * @param userId 用户ID
      * @param qrCodeId 二维码唯一id
      * @param SDKID
      */
      aphead.PCLoginQrCodeResult scanPCLoginQrCode(1:string elogo,2:i64 userId,3:string qrCodeId,4:i64 SDKID);

      /**
        * 移动端扫PC端登录二维码确认接口
        * @param elogo 服务器elogo
        * @param userId 用户ID
        * @param qrCodeId 二维码唯一id
        * @param SDKID
        */
        aphead.PCLoginQrCodeResult scanPCLoginQrCodeConfirm(1:i64 userId,2:string qrCodeId,3:byte confirmLogin,4:i64 SDKID);

        /**
        * 获取二维码状态接口
        * @param qrCodeId 二维码唯一id
        * @param qrCodeIdEv 二维码唯一id加密值
        * @param SDKID
        */
        aphead.PCLoginQrCodeResult getPCLoginQrCodeStatus(1:string qrCodeId,2: string qrCodeIdEv,3:i64 SDKID);
      /****************************************PC扫码登录 end*********************************/

       /**
       * 根据用户ID集合查询用户的设备信息,豆豆版本
       * @param userIds 用户ID集合
       * @param SDKID SDKID
       */
       aphead.ResultLoginRecord getDeviceInfoListByUserIds(1:list<i64> userIds,2:i64 SDKID);

       /**
        * 根据时间戳批量删除用户登录记录
        */
        aphead.Result deleteUserRecordsByDate(1:i64 timestamp);

	/****************************************PC历史账户一键登录 start*********************************/
	/**
	* 校验最后一次PC扫码登录设备接口
	* @param userId 用户ID
	* @param macCode macCode
	* @param SDKID
	*/
	aphead.PCLoginQrCodeResult checkRecentQcLoginDevice(1:i64 userId, 2:string macCode,3:i64 SDKID);

	/**
	* PC端快速登录接口
	* @param userId 用户ID
	* @param loginData PC端登录信息
	* @param SDKID
	*/
	aphead.PCLoginQrCodeResult pcQuickLogin(1:i64 userId, 2:string loginData,3:i64 SDKID);

	/****************************************PC历史账户一键登录 end*********************************/

	/****************************************后台双因子登录 start*********************************/

	/**
    * 推送确认登录操作至客户端的接口
    * @param userId 用户ID
    * @param bifactorLogin 后台双因子认证实体
    * @param SDKID
    */
    aphead.ResultBifactor pushBifactorLogin(1:i64 userId, 2:aphead.BifactorLogin bifactorLogin,3:i64 SDKID);

    /**
    * 客户端确认/拒绝登录接口
    * @param userId 用户ID
    * @param bifactorLogin 后台双因子认证实体
    * @param SDKID
    */
    aphead.ResultBifactor confirmLogin(1:i64 userId, 2:aphead.BifactorLogin bifactorLogin,3:i64 SDKID);

    /**
    * 查询用户是否确认登录状态接口
    * @param userId 用户ID
    * @param bifactorLogin 后台双因子认证实体
    * @param SDKID
    */
    aphead.ResultBifactor pollingBifactorLogin(1:i64 userId, 2:aphead.BifactorLogin bifactorLogin,3:i64 SDKID);

	/****************************************后台双因子登录 end***********************************/

    /**
    * 根据登录时间和用户ID集合批量查询登录记录
    * @param userIds 用户ID集合
    * @param startTime 开始时间
    * @param endTime 结束时间
    * @param SDKID
    */
      list<aphead.LoginRecord> getLoginRecordListByUserIdsAndLoginTime(1:list<i64> userIds,2:i64 startTime,3:i64 endTime,4:i64 SDKID);
	  
	/**
	* 根据userID查询的符合条件的集合对象
	*/
	list<aphead.UserBindingDeviceInfo> queryUserBindingDeviceInfo(1: i64 userID,2:i64 SDKID);
	
	/**
	* 根据USERID,SDKID,KEY删除用户绑定的设备信息
	*/
	aphead.Result deleteUserBindDeviceInfo(1: i64 userID,2:i64 SDKID,3:string key);
	
	/**
	* 更新用户绑定的设备信息
	*/
	i64 editUserBindDeviceInfo(1: aphead.UserBindingDeviceInfo userBindingDeviceInfo);
	
	/**
	* 添加用户绑定的设备信息
	*/
	i64 addUserBindDeviceInfo(1: aphead.UserBindingDeviceInfo userBindingDeviceInfo)

}

`

func TestThrift(t *testing.T) {
	tree, err := Parse(thriftCode)
	if err != nil {
		fmt.Println(err.Error())
	}
	node.OutTree(thriftCode, tree)

}
