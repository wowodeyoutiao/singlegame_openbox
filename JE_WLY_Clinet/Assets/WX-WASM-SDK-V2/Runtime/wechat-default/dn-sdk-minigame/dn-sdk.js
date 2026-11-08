import { SDK }from './index.js';

SDK.setDebug(false);

const sdk = new SDK({
    //数据源ID，数字，必填。
    user_action_set_id: 1213615927,
    //加密key，必填。
    secret_key: '1652d1066ffe1393acc9045b9eb550a9',
    //微信小游戏AppID，wx开头，必填。
    appid: 'wx45449e394d2f152b'
});

GameGlobal.setOpenId = (args) => {
    //设置opneid，必须先设置openid再上报注册行为。setopenId为同步方法，设需完请交即上报注册行为。
    sdk.setOpenId(args.openId);
    console.log("设置opneid: "+args.openId);
};

GameGlobal.setUnionId = (args) => {
    //设置opneid，必须先设置openid再上报注册行为。setopenId为同步方法，设需完请交即上报注册行为。
    sdk.setUnionId(args.unionid);
    console.log("设置opneid: "+args.unionid);
};

GameGlobal.onRegister = () => {
    sdk.onRegister();
    console.log("用户注册成功！");
}

GameGlobal.createRole = (args) => {
    //事件注册在GameGloble里
    console.log("创建角色: "+args.userName);
    sdk.onCreateRole(args.userName);
};

GameGlobal.purchase = (args) => {
    console.log("支付价格: " + args.price + "分");
    // price 单位：分
    sdk.onPurchase(args.price);
};

GameGlobal.tutorialFinish = () => {
    sdk.onTutorialFinish()
    console.log("完成新手教程!");
}