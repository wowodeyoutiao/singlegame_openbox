using System;
using System.Collections;
using UnityEngine;
using UnityEngine.UI;

public class BQNoPwdLoginCtr: BQSDKUILogicBase
{
    public InputField iptPhone;
    public InputField iptCode;

    public Button btnClose;

    public Button btnUserProtocol;
    public Button btnPriviteProtocol;
    public Button btnProtectChildrenProtocol;

    public Toggle tgAgreeProtocol;

    public Button btnSendCodeToPhone;
    public Button btnLogin;
    public GameObject goDisableLogin;
    public Text txtCodeTips;

    bool agree = false;
    public override void OnShow()
    {
        base.OnShow();

        agree = tgAgreeProtocol.isOn;
    }
    protected override void AddEventListener()
    {
        btnClose?.onClick.AddListener(() => { this.ParentCtr.GoToLogin(); });

        btnUserProtocol?.onClick.AddListener(() => { this.ParentCtr.OnShowUserProtocol(); });
        btnPriviteProtocol?.onClick.AddListener(() => { this.ParentCtr.OnShowPriviteProtocol(); });
        btnProtectChildrenProtocol?.onClick.AddListener(() => { this.ParentCtr.OnShowProtectChildrenProtocol(); });

        tgAgreeProtocol?.onValueChanged.AddListener((bool ison) => { agree = ison; });

        btnSendCodeToPhone?.onClick.AddListener(OnSendCodeToPhone);
        btnLogin?.onClick.AddListener(OnLogin);
    }

    protected override void RemoveEventListener()
    {
        btnClose?.onClick.RemoveAllListeners();

        btnUserProtocol?.onClick.RemoveAllListeners();
        btnPriviteProtocol?.onClick.RemoveAllListeners();
        btnProtectChildrenProtocol?.onClick.RemoveAllListeners();

        tgAgreeProtocol?.onValueChanged.RemoveAllListeners();

        btnSendCodeToPhone?.onClick.RemoveAllListeners();
        btnLogin?.onClick.RemoveAllListeners();
    }

    public override void OnHide()
    {
        base.OnHide();
    }

    #region UI Button Event
    private void OnSendCodeToPhone()
    {
        
    }

    private void OnLogin()
    {
        
    }
    #endregion
}