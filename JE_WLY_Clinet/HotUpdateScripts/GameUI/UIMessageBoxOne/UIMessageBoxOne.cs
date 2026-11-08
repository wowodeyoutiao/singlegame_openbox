using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using System;
using JEngine.Core;

public class UIMessageBoxOne : UIBase<UIModelMessageBoxOne, UIViewMessageBoxOne>
{
    public class CallBack 
    {
        public string showStr;
        public string btnStr;
        public int cType;
        public int timeCount;
        public Action endFun;
    }

    public override UILayer GetLayer()
    {
        return UILayer.Top;
    }
    private CallBack _callBack ;
    public override void OnInit()
    {
        GetView().Btn_Yes.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            CloseClick();
        });
    }

    private void CloseClick()
    {
        ClearJa();
        UITool.Close<UIMessageBoxOne>();
        _callBack.endFun?.Invoke();
    }

    private TimeItem Ja;

    private void ClearJa()
    {
        if (Ja != null)
        {
            GameHelper.Instance.RemoveInvoke(Ja);
            Ja = null;
        }
    }
    public override void OnEnter(params object[] args)
    {
        _callBack = args[0] as CallBack;
        
        if (_callBack.cType == 0)
            GetView().Text_des.text = _callBack.showStr;
        else
        {
            Ja = GameHelper.Instance.AddInvoke(new TimeItem(_callBack.timeCount, 1, () => {
                GetView().Text_des.text = String.Format(_callBack.showStr, _callBack.timeCount);
                _callBack.timeCount--;
                if (_callBack.timeCount == 0)
                    CloseClick();
            }));
        }

        GetView().Text_Yes.text = string.IsNullOrEmpty(_callBack.btnStr) ? "确定" : _callBack.btnStr;
    }

    public override void OnExit()
    {
        ClearJa();
    }

    public override void OnDestroy()
    {
        ClearJa();
    }
}