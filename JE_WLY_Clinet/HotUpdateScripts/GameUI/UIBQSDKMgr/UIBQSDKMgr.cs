
using HotUpdateScripts;
using MetaFramework.UI;

using System;
using System.Collections.Generic;
using System.Security.Principal;
using System.Text.RegularExpressions;
using UnityEngine;
using UnityEngine.UI;

public class UIModelBQSDKMgr : IUIModel
{

}

public class UIViewBQSDKMgr : MonoBehaviour, IUIView
{
    public BQLoginCtr loginCtr;
    public BQBindPhoneCtr bindPhoneCtr;
    public BQForgetPwdCtr forgetPwdCtr;
    public BQNoPwdLoginCtr noPwdLoginCtr;
    public BQRegisterCtr registerCtr;
    public BQCertification certificationCtr;

    public void Awake()
    {
    }

    public void Init(GameObject go)
    {
    }
    private void OnDestroy()
    {
        UITool.RemoveUI<UIBQSDKMgr>();
    }
}

public class UIBQSDKMgr : UIBase<UIModelBQSDKMgr, UIViewBQSDKMgr>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    
    public override void OnInit()
    {
        GetView().loginCtr?.Init(this);
        GetView().bindPhoneCtr?.Init(this);
        GetView().forgetPwdCtr?.Init(this);
        GetView().noPwdLoginCtr?.Init(this);
        GetView().registerCtr?.Init(this);
        GetView().certificationCtr?.Init(this);
    }

    BQSDKUILogicBase curCtrl = null;

    public override void OnEnter(params object[] args)
    {
        ChangeUI(GetView().loginCtr);
    }

    public override void OnExit()
    {
        curCtrl?.OnHide();

        GetView().loginCtr?.Exit();
        GetView().bindPhoneCtr?.Exit();
        GetView().forgetPwdCtr?.Exit();
        GetView().noPwdLoginCtr?.Exit();
        GetView().registerCtr?.Exit();
        GetView().certificationCtr?.Exit();
    }

    #region UI Jumping
    public void GoToLogin()
    {
        ChangeUI(GetView().loginCtr);
    }

    public void GoToBindPhoneOrUnBind()
    {
        ChangeUI(GetView().bindPhoneCtr);
    }

    public void GoToForgetPwd()
    {
        ChangeUI(GetView().forgetPwdCtr);
    }

    public void GoToNoPwdLogin()
    {
        ChangeUI(GetView().noPwdLoginCtr);
    }

    public void GoToRegister()
    {
        ChangeUI(GetView().registerCtr);
    }

    public void GoToCertification()
    {
        ChangeUI(GetView().registerCtr);
    }

    private void ChangeUI(BQSDKUILogicBase newCtrl)
    {
        curCtrl?.OnHide();
        curCtrl = newCtrl;
        curCtrl.OnShow();
    }

    #endregion

    #region Show Web
    public void OnShowUserProtocol()
    {
        UITool.ShowTips("功能完善中！");
    }

    public void OnShowPriviteProtocol()
    {
        UITool.ShowTips("功能完善中！");
    }

    public void OnShowProtectChildrenProtocol()
    {
        UITool.ShowTips("功能完善中！");
    }
    #endregion

    #region 

    public bool IsChooseProtocol(bool agree)
    {
        if (!agree)
        {
            UITool.ShowTips("请先勾选同意协议");
            return false;
        }

        return true;
    }

    public bool CheckStrIsNullOrEmpty(string str, string tipsWhenStrIsNull = null)
    {
        bool result = string.IsNullOrEmpty(str);
        if (result && !string.IsNullOrEmpty(tipsWhenStrIsNull))
        {
            UITool.ShowTips(tipsWhenStrIsNull);
        }

        return result;
    }

    public bool ComparePwdAdnEnterPwd(string pwd, string enterPwd, string tipsWhenCompareFaild = null)
    { 
        bool result = pwd.Equals(enterPwd);

        if (!result && !string.IsNullOrEmpty(tipsWhenCompareFaild))
        {
            UITool.ShowTips(tipsWhenCompareFaild);
        }

        return result;
    }

    public bool AccountFormatIsPass(string account)
    {
        // 长度判断（一个中文算一个字符）
        if (account.Length > 11)
        {
            UITool.ShowTips("账号长度不能超过11个字符");
            return false;
        }

        // 只允许中文、英文字母、数字
        if (!Regex.IsMatch(account, @"^[\u4e00-\u9fa5a-zA-Z0-9]+$"))
        {
            UITool.ShowTips("账号应是文字、字母、数字的组合");
            return false;
        }
            

        // 检查至少包含两种类型
        bool hasChinese = Regex.IsMatch(account, @"[\u4e00-\u9fa5]");
        bool hasLetter = Regex.IsMatch(account, @"[a-zA-Z]");
        bool hasDigit = Regex.IsMatch(account, @"\d");

        int typeCount = 0;
        if (hasChinese) typeCount++;
        if (hasLetter) typeCount++;
        if (hasDigit) typeCount++;

        if (typeCount < 2)
        {
            UITool.ShowTips("账号应是文字、字母、数字的组合");
            return false;
        }

        return true;
    }

    public bool PwdFormatIsPass(string pwd)
    {
        if (pwd.Length > 8)
        {
            UITool.ShowTips("密码长度不能超过8个字符");
            return false;
        }

        // 只允许英文字母、数字
        if (!Regex.IsMatch(pwd, @"^[a-zA-Z0-9]+$"))
        {
            UITool.ShowTips("密码应是字母、数字的组合");
            return false;
        }


        // 检查必须包含两种类型
        bool hasLetter = Regex.IsMatch(pwd, @"[a-zA-Z]");
        bool hasDigit = Regex.IsMatch(pwd, @"\d");

        int typeCount = 0;
        if (hasLetter) typeCount++;
        if (hasDigit) typeCount++;

        if (typeCount < 2)
        {
            UITool.ShowTips("密码应是字母、数字的组合");
            return false;
        }

        return true;
    }

    #endregion
}

