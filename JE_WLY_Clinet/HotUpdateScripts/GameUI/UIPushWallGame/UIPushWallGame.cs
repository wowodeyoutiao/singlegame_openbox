using MetaFramework.UI;
using System;
using System.Collections;
using UnityEngine;

public class UIModelPushWallGame : IUIModel
{

}

public class UIViewPushWallGame : MonoBehaviour, IUIView
{
    public void Awake()
    {
    }

    public void Init(GameObject go)
    {
    }
    private void OnDestroy()
    {
        UITool.RemoveUI<UIPushWallGame>();
    }
}

public class UIPushWallGame : UIBase<UIModelPushWallGame, UIViewPushWallGame>
{
    //public 
    public ComboHandle combHandle;//手柄管理类
    public ComboItemsMgr comboItemsMgr;//combo管理类

    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    public override void OnInit()
    {
        combHandle.Init(ComboHandleClickCB);
        comboItemsMgr.Init(ComboResultCB);
    }

    public override void OnEnter(params object[] args)
    {
        
    }

    public override void OnExit()
    {
        
    }

    private void ComboHandleClickCB(int comboNo)
    {
        if(!comboItemsMgr.IsCanInput()) return; //非接收输入状态

        comboItemsMgr?.Match(comboNo); // 等待 ComboResultCB
    }

    //ComboHandleClickCB负责传递操作，这里等待combo结果
    private void ComboResultCB(bool isSuccess)
    {
        if (isSuccess)
        {
            //成功后发送效果增益
            //等待5秒进入下一轮技能选择
        }
        else
        {
            //失败后 等待5秒进入下一轮技能选择
        }
    }
}