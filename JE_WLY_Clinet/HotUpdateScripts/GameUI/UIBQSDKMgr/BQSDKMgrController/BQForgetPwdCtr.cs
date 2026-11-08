using System.Collections;
using UnityEngine;
using UnityEngine.UI;

public class BQForgetPwdCtr : BQSDKUILogicBase
{
    public Button btnClose;

    //确认账号
    public InputField iptAccount;
    public Button btnEnterAccount;

    //输入手机号并发送验证码
    public InputField iptPhoneNum;
    public InputField iptCode;
    public Text txtCodeTips;
    public Button btnSendCodeToPhone;
    public Button btnEnterPNAndCode;
    public GameObject goEnable;

    //确认密码
    public Text txtPwd;
    public Button btnCopyPwd;

    public override void OnShow()
    {
        base.OnShow();
    }
    protected override void AddEventListener()
    {
        btnClose?.onClick.AddListener(() => { this.ParentCtr.GoToLogin(); });

    }

    protected override void RemoveEventListener()
    {

    }

    public override void OnHide()
    {
        base.OnHide();
    }
}