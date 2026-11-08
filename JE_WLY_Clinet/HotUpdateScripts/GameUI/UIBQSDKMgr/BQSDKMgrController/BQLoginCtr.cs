using MetaFramework.UI;
using System;
using System.Collections;
using UnityEngine;
using UnityEngine.UI;

public class BQLoginCtr : BQSDKUILogicBase
{
    public Button btnLogin;

    public Button btnUserProtocol;
    public Button btnPriviteProtocol;
    public Button btnProtectChildrenProtocol;

    public Button btnNoPwdLoginEntrance;
    public Button btnForgetPwdEntrance;
    public Button btnRegister;

    public Toggle tgAgreeProtocol;

    bool agree = false;
    public override void OnShow()
    {
        base.OnShow();

        agree = tgAgreeProtocol.isOn;
    }
    protected override void AddEventListener()
    {
        btnLogin?.onClick.AddListener(OnLoginClick);

        btnUserProtocol?.onClick.AddListener(() => { this.ParentCtr.OnShowUserProtocol(); });
        btnPriviteProtocol?.onClick.AddListener(() => { this.ParentCtr.OnShowPriviteProtocol(); });
        btnProtectChildrenProtocol?.onClick.AddListener(() => { this.ParentCtr.OnShowProtectChildrenProtocol(); });

        btnNoPwdLoginEntrance?.onClick.AddListener(() => { this.ParentCtr.GoToNoPwdLogin(); });
        btnForgetPwdEntrance?.onClick.AddListener(() => { this.ParentCtr.GoToForgetPwd(); });
        btnRegister?.onClick.AddListener(() => { this.ParentCtr.GoToRegister(); });

        tgAgreeProtocol?.onValueChanged.AddListener((bool ison) => { agree = ison; });
    }

    protected override void RemoveEventListener()
    {
        btnLogin?.onClick.AddListener(OnLoginClick);

        btnUserProtocol?.onClick.RemoveAllListeners();
        btnPriviteProtocol?.onClick.RemoveAllListeners();
        btnProtectChildrenProtocol?.onClick.RemoveAllListeners();

        btnNoPwdLoginEntrance?.onClick.RemoveAllListeners();
        btnForgetPwdEntrance?.onClick.RemoveAllListeners();
        btnRegister?.onClick.RemoveAllListeners();

        tgAgreeProtocol?.onValueChanged.RemoveAllListeners();
    }

    public override void OnHide()
    {
        base.OnHide();
    }


    #region UI Button Event
    private void OnLoginClick()
    {
        if (!agree)
        {
            UITool.ShowTips("请先勾选同意协议");
            return;
        }
    }

    #endregion
}