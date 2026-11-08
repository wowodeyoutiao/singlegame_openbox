mergeInto(LibraryManager.library,{
    SetOpenId:function(openId){
        GameGlobal.setOpenId({
            openId: UTF8ToString(openId)
        });
    },

    SetUnionId:function(unionid){
        GameGlobal.setUnionId({
            unionid: UTF8ToString(unionid)
        });
    },
    
    OnRegister:function(){
        GameGlobal.onRegister();
    },
    
    CreateRole:function(userName){
        GameGlobal.createRole({
            userName: UTF8ToString(userName)
        });
    },
    
    Purchase:function(price){
        GameGlobal.purchase({
            price: price
        });
    },
    
    TutorialFinish:function(){
        GameGlobal.tutorialFinish();
    },
});