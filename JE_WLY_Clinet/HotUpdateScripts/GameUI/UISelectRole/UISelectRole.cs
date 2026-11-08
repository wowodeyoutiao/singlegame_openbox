using Cysharp.Threading.Tasks;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.LitJson;
using JEngine.MyNet;
using JEngine.UI;
using MetaFramework.UI;
using Spine.Unity;
using System;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class UISelectRole : UIBase<UIModelSelectRole, UIViewSelectRole>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    private int Select_Idx = -1;
    public override void OnInit()
    {
        GetView().Btn_Creat.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnCreatBtnClick();
        });

        GetView().Btn_Delete.onClick.AddListener(OnDeleteBtnClick);

        GetView().Btn_StarGame.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnStarGameBtnClick();
        });

        GetView().btn_Back.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UISelectRole>();
            TCPClientMgr.Instance.ClearIpPool();
        });

        GetView().Btn_Delete.gameObject.SetActive(InitJEngine.Instance.IsFreeLogin);

    }


    private void OnCreatBtnClick()
    {
        UITool.Open<UICreatRole>();
        // UITool.Close<UISelectRole>();
    }

    private void OnDeleteBtnClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (Select_Idx > -1)
        {
            long deleteUid = RoleM.Instance.LoginCharacterInfo[Select_Idx].uid;
            string selectInfo = JsonMapper.ToJson(new HTTP_LoginrateRequest()
            {
                uid = deleteUid
            });
            Action yesFun = () =>
            {
                HttpRequestHelper.SendPost<HTTP_LoginrateResponse>(GameM.Instance.HttpUrl + GameM.Instance.Role_DeleteUrl, NetworkHelper.HTTPMethod.DELETE, RoleM.Instance.Login_Token, selectInfo,
                (response) =>
                {
                    if (response.code == 0)
                    {
                        int dIdx = RoleM.Instance.LoginCharacterInfo.FindIndex(x => x.uid == deleteUid);
                        RoleM.Instance.LoginCharacterInfo.RemoveAt(dIdx);
                        if (RoleM.Instance.LoginCharacterInfo.Count == 0)
                        {
                            OnCreatBtnClick();
                            UITool.Close<UISelectRole>();
                        }
                        else
                        {
                            RefreshRoles();
                            SetBtnState();
                        }
                    }
                });
            };
            yesFun();
            //UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack { 
            //    ShowStr = $"是否确定删除当前角色:{RoleM.Instance.LoginCharacterInfo[Select_Idx].name}?",
            //    YesFun = () => {
            //        UITool.Open<UIMessageDestoryRole>(null, yesFun);
            //    },
            //});
        }
    }

    private void RefreshRoles()
    {
        var RoleInfo = RoleM.Instance.LoginCharacterInfo;
        for (int i = 0; i < 2; i++)
        {
            GameObject Item = GetView().Root_Btn.transform.Find("itemRole" + i).gameObject;
            Item.SetActive(i < RoleInfo.Count);
            if (i < RoleInfo.Count)
            {
                Item.transform.Find("bg/Text").GetComponent<Text>().text = RoleInfo[i].name;
                Item.transform.Find("bg/txt_level").GetComponent<Text>().text = RoleInfo[i].level + "级";

                Item.transform.Find("on/Text").GetComponent<Text>().text = RoleInfo[i].name;
                Item.transform.Find("on/txt_level").GetComponent<Text>().text = RoleInfo[i].level + "级";

                var iconStr = "img_casting_tx_" + ((RoleInfo[i].job - 1) * 2 + RoleInfo[i].sex);
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CreatRoleAtlas, iconStr, JUIExtensions.img(Item.transform, "bg/icon"));
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CreatRoleAtlas, iconStr, JUIExtensions.img(Item.transform, "on/icon"));
                Item.GetComponent<Button>().onClick.RemoveAllListeners();
                var idx = i;
                Item.GetComponent<Button>().onClick.AddListener(() =>
                {
                    OnSelectRoleBtnClick(idx);
                });
            }

        }
        OnSelectRoleBtnClick(0);

    }

    public override void OnLayerChange()
    {
        GameUtil.Instance.SetEffectLayer(GetView().img_job.transform, _canvasOrder + 1);
    }

    //开始游戏
    private void OnStarGameBtnClick()
    {
        UITool.Open<UILoading>((ui) =>
        {
            HttpC.OnSelectRoleToLogin(Select_Idx);
        }, new UILoading.Loading { _time = 0f, _str = "登录..." });

        string selectedServerName = GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx].name;
        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties["page"] = "登录界面";
        properties["component"] = "选角界面";
        properties["position"] = "";
        properties["event_click"] = $"进入游戏(服务器：{selectedServerName})";
        TDAnalyticsUtil.Track("click_event", properties, true);
    }

    public static string[] JobIconName = { "", "icon_1", "icon_2", "icon_3" };
    private void OnSelectRoleBtnClick(int idx)
    {
        for (int i = 0; i < 2; i++)
        {
            //GetView().Root_Btn.transform.Find("itemRole" + i).transform.localScale = idx == i ? Vector3.one * 1.2f : Vector3.one;
            GetView().Root_Btn.transform.Find("itemRole" + i + "/on").gameObject.SetActive(idx == i);
        }

        Select_Idx = idx;
        int job = RoleM.Instance.LoginCharacterInfo[Select_Idx].job;
        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CreatRoleAtlas, JobIconName[job], GetView().img_job);
        GetView().txt_JobName.text = CFG_Item.GetJobStrById(job);

        var strs = UIModelCreatRole.JobDesDic[(JobType)job];
        for (int i = 0; i < strs.Length; i++)
        {
            GetView().JobDesTxts[i].text = strs[i];
        }
        SetRolec((JobType)job, (SexType)RoleM.Instance.LoginCharacterInfo[Select_Idx].sex);

        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties["page"] = "登录界面";
        properties["component"] = "选角界面";
        properties["position"] = "";
        properties["event_click"] = $"选择角色(角色{idx + 1})";
        TDAnalyticsUtil.Track("click_event", properties, true);
    }

    private void SetRolec(JobType jobType, SexType sexType)
    {
        if (GetView().root_Role.childCount > 0)
        {
            for (int i = 0; i < GetView().root_Role.childCount; i++)
            {
                PoolHelper.Instance.RecycObj(GetView().root_Role.GetChild(i).gameObject);
            }
        }

        GameHelper.Instance.SetRoleAsync(jobType, sexType, (_mainRole) =>
        {
            _mainRole.transform.SetParent(GetView().root_Role);
            _mainRole.transform.localScale = Vector3.one;
            _mainRole.transform.localPosition = Vector3.zero;
            _mainRole.SetActive(true);

            SpineUtil.Instance.PlayQueue(_mainRole.transform.Find("spine").GetComponent<SkeletonGraphic>(), new List<string> { "skill2", "skill3", "stand" }, null, true);
        }).Forget();

    }
    public override void OnEnter(params object[] args)
    {
        RefreshRoles();
        SetBtnState();
    }

    private void SetBtnState()
    {
        GetView().Btn_Creat.gameObject.SetActive(CFG_XMLConfig.xmlConfig.BaseValue.global.rolenum > RoleM.Instance.LoginCharacterInfo.Count);
    }

    public override void OnExit()
    {

    }

}