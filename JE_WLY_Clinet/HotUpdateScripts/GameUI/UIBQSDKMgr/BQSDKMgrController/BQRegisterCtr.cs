using MetaFramework.UI;
using System;
using System.Collections;
using UnityEngine;
using UnityEngine.UI;

public class BQRegisterCtr : BQSDKUILogicBase
{
    public Button btnClose;
    public InputField iptAccount;
    public InputField iptPassword;
    public InputField iptPasswordAgain;

    public Button btnUserProtocol;
    public Button btnPriviteProtocol;
    public Button btnProtectChildrenProtocol;

    public UISwitchButton swbPwd;
    public UISwitchButton swbEnterPwd;

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
        btnClose?.onClick.AddListener(() => { this.ParentCtr.GoToLogin(); });

        btnUserProtocol?.onClick.AddListener(() => { this.ParentCtr.OnShowUserProtocol(); });
        btnPriviteProtocol?.onClick.AddListener(() => { this.ParentCtr.OnShowPriviteProtocol(); });
        btnProtectChildrenProtocol?.onClick.AddListener(() => { this.ParentCtr.OnShowProtectChildrenProtocol(); });

        btnRegister?.onClick.AddListener(OnRegisterClick);

        tgAgreeProtocol?.onValueChanged.AddListener((bool ison) => { agree = ison; });

        swbPwd?.Init(true, (bool isOn) => { SetInputFiledContentType(iptPassword, isOn); });
        swbEnterPwd?.Init(true, (bool isOn) => { SetInputFiledContentType(iptPasswordAgain, isOn); });
    }

    protected override void RemoveEventListener()
    {
        btnClose?.onClick.RemoveAllListeners();

        btnUserProtocol?.onClick.RemoveAllListeners();
        btnPriviteProtocol?.onClick.RemoveAllListeners();
        btnProtectChildrenProtocol?.onClick.RemoveAllListeners();

        btnRegister?.onClick.RemoveAllListeners();

        tgAgreeProtocol?.onValueChanged.RemoveAllListeners();
    }

    public override void OnHide()
    {
        base.OnHide();
    }

    private void SetInputFiledContentType(InputField input,bool isShowPwd)
    {
        input.contentType = isShowPwd ? InputField.ContentType.Password : InputField.ContentType.Standard;
        input.ForceLabelUpdate();
    }

    #region UI Button Event
    private void OnRegisterClick()
    {
        if (!this.ParentCtr.IsChooseProtocol(agree)) return;

        if (this.ParentCtr.CheckStrIsNullOrEmpty(iptAccount.text,"账号不能为空！")) return;

        if (!this.ParentCtr.AccountFormatIsPass(iptAccount.text)) return;//账号格式不正确

        if (this.ParentCtr.CheckStrIsNullOrEmpty(iptPassword.text, "密码不能为空！")) return;

        if (!this.ParentCtr.PwdFormatIsPass(iptPassword.text)) return; //密码格式不正确

        if (this.ParentCtr.CheckStrIsNullOrEmpty(iptPasswordAgain.text, "确认密码不能为空！")) return;

        if (this.ParentCtr.ComparePwdAdnEnterPwd(iptPassword.text, iptPasswordAgain.text, "密码与确认密码不一致！请确认。")) return;


    }
    #endregion
}