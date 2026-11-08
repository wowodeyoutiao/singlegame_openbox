using System.Collections;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.LitJson;
using JEngine.Core;
using HotUpdateScripts;
using Spine.Unity;
using System.Collections.Generic;
using System.Linq;
using Cysharp.Threading.Tasks;
using JEngine.UI;

public class UICreatRole : UIBase<UIModelCreatRole, UIViewCreatRole>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    private HTTP_CreateCharacterRequest CreatInfo = new HTTP_CreateCharacterRequest
    {
        sex = (int)SexType.SexGril,
        clickid = GameUtil.Instance.GetWXLaunchParamClickId()
    };

    public override void OnInit()
    {
        GetView().btn_Creat.onClick.AddListener(OnCreatBtnClick);
        GetView().btn_randomName.onClick.AddListener(delegate { 
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnRandmNameBtnClick();
        });

        for (int i = 0; i < GetView().JobSelectBtns.Length; i++)
        {
            int idx = i ;
            GetView().JobSelectBtns[i].onClick.AddListener(delegate {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnJobChoiscClick(idx); 
            });
        }

        //GetView().btn_Gril.onClick.AddListener(() => {
        //    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        //    OnSexChoiscClick(SexType.SexGril);
        //});

        //GetView().btn_Boy.onClick.AddListener(() => {
        //    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        //    OnSexChoiscClick(SexType.SexBoy);
        //});

        GetView().btn_back.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UICreatRole>();
        });
    }

    private void SetSex(SexType sexIDx)
    {
        if (IsClick)
        {
            UITool.ShowTips("操作频繁，请稍后");
            return;
        }
        if (!IsClick)
            IsClick = true;

        JUIExtensions.SetGray(GetView().img_Boy, sexIDx != SexType.SexBoy);
        JUIExtensions.SetGray(GetView().img_Gril, sexIDx != SexType.SexGril);
        CreatInfo.sex = (int)sexIDx;

        SetRole((JobType)CreatInfo.job, sexIDx);

        string sex = sexIDx switch
        {
            SexType.SexBoy => "男",
            SexType.SexGril => "女",
            _ => "无",
        };

        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties["page"] = "登录界面";
        properties["component"] = "创角界面";
        properties["position"] = "";
        properties["event_click"] = $"选择性别({sex})";
        TDAnalyticsUtil.Track("click_event", properties, true);
    }

    bool IsClick = false;

    private void OnJobChoiscClick(int idx)
    {
        if (IsClick)
        {
            UITool.ShowTips("操作频繁，请稍后");
            return ;
        }

        JobType jobIdx = (JobType)(idx % 3 + 1);

        GetView().root_Select.localPosition = GetView().JobSelectBtns[idx].transform.localPosition;

        CreatInfo.job = (int)jobIdx;

        GetView().txt_Job.text = CFG_Item.GetJobStrById((int)jobIdx);
        var strs = UIModelCreatRole.JobDesDic[jobIdx];
        for (int i = 0; i < strs.Length; i++)
        {
            GetView().JobDesTxts[i].text = strs[i];
        }

        var sexIDx = idx > 2 ? SexType.SexGril : SexType.SexBoy;
        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CreatRoleAtlas, UISelectRole.JobIconName[(int)jobIdx], GetView().img_job);
        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CreatRoleAtlas, "img_casting_lh_"+ ((CreatInfo.job -1) * 2 + (int)sexIDx), GetView().img_Role);

        SetSex(sexIDx);

        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties["page"] = "登录界面";
        properties["component"] = "创角界面";
        properties["position"] = "";
        properties["event_click"] = $"选择职业({GetView().txt_Job.text})";
        TDAnalyticsUtil.Track("click_event", properties, true);
    }

    private Dictionary<string, string[]> RoleBronAnim = new Dictionary<string, string[]>
    {
        [$"{JobType.Warrior}{SexType.SexBoy}"] = new string[] {GameConst.ModelAnim.skill2.ToString(), GameConst.ModelAnim.skill3.ToString(),GameConst.ModelAnim.attack.ToString(), GameConst.ModelAnim.stand.ToString()},
        [$"{JobType.Warrior}{SexType.SexGril}"] = new string[] {GameConst.ModelAnim.skill2.ToString(), GameConst.ModelAnim.skill3.ToString(), GameConst.ModelAnim.skill1.ToString(), GameConst.ModelAnim.stand.ToString()},
        [$"{JobType.Mage}{SexType.SexBoy}"] = new string[] {GameConst.ModelAnim.skill1.ToString(), GameConst.ModelAnim.skill3.ToString(), GameConst.ModelAnim.stand.ToString()},
        [$"{JobType.Mage}{SexType.SexGril}"] = new string[] {GameConst.ModelAnim.skill1.ToString(), GameConst.ModelAnim.skill3.ToString(), GameConst.ModelAnim.stand.ToString()},
        [$"{JobType.Sacrifice}{SexType.SexBoy}"] = new string[] {GameConst.ModelAnim.skill1.ToString(), GameConst.ModelAnim.attack.ToString(),  GameConst.ModelAnim.stand.ToString()},
        [$"{JobType.Sacrifice}{SexType.SexGril}"] = new string[] {GameConst.ModelAnim.skill1.ToString(), GameConst.ModelAnim.attack.ToString(), GameConst.ModelAnim.stand.ToString()},

    };

    private void SetRole(JobType jobType, SexType sexType)
    {
        int Count = GetView().root_role.childCount;
        if (Count > 0)
        {
            for (int i = 0; i < Count; i++)
            {
                PoolHelper.Instance.RecycObj(GetView().root_role.GetChild(0).gameObject);
            }
        }
        
        GameObject Eff = null;
        GameHelper.Instance.SetRoleAsync(jobType, sexType, (_mainRole) =>
        {
           
            _mainRole.transform.SetParent(GetView().root_role);
            _mainRole.transform.localScale = Vector3.one;
            _mainRole.transform.localPosition = Vector3.zero;
            _mainRole.SetActive(true);

            string booEff = $"{jobType.ToString()}{sexType.ToString()}_Bron";
            string booEffPath = GameConst.GetBornPatchByJobAndSex(jobType, sexType);
            if (!string.IsNullOrEmpty(booEffPath))
            {
                _ = GameHelper.Instance.SetAsyncObj(booEff, booEffPath, (go) =>
                {
                    go.SetActive(true);
                    go.transform.SetParent(GetView().root_role);
                    go.transform.localPosition = Vector3.zero;
                    go.transform.localScale = Vector3.one;
                    IsClick = false;
                    Eff = go;
                },true);
            }

            SpineUtil.Instance.PlayQueue(_mainRole.transform.Find("spine").GetComponent<SkeletonGraphic>(), RoleBronAnim[$"{jobType}{sexType}"].ToList(), () => {
                if (Eff != null)
                {
                    PoolHelper.Instance.RecycObj(Eff);
                    Eff = null;
                }
            }, true);
        }).Forget();
            
        //GameHelper.Instance.SetDeayTimeDo(() => { IsClick = false; },1f);
    }

    private void OnSexChoiscClick(SexType SexIdx)
    {
        CreatInfo.sex = (int)SexIdx;
        SetSex((SexType) SexIdx);

        string sex = SexIdx switch
        {
            SexType.SexBoy => "男",
            SexType.SexGril => "女",
            _ => "无",
        };

        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties["page"] = "登录界面";
        properties["component"] = "创角界面";
        properties["position"] = "";
        properties["event_click"] = $"选择性别({sex})";
        TDAnalyticsUtil.Track("click_event", properties, true);
    }
    private void OnCreatBtnClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (GetView().input_Name.text == "")
        {
            UITool.ShowTips("名字为空");
            return;
        }
        CreatInfo.name = GetView().input_Name.text;

        HttpRequestHelper.SendPost<HTTP_CreateCharacterResponse>(GameM.Instance.HttpUrl + GameM.Instance.Role_CreatUrl, NetworkHelper.HTTPMethod.PUT, RoleM.Instance.Login_Token, JsonMapper.ToJson(CreatInfo),
            (_creat) => {
                if (_creat.code == 0)
                {
                    SDKCtr.Instance.DNCreateRole(_creat.characterList.name);

                    if (!isNullGo)
                    {
                        RoleM.Instance.LoginCharacterInfo.Insert(0, _creat.characterList);
                        //UITool.Open<UISelectRole>(null,true);

                        UITool.Open<UILoading>((ui) => {
                            HttpC.OnSelectRoleToLogin(0, () =>
                            {
                                UITool.Close<UICreatRole>();
                                UITool.Close<UISelectRole>();
                            });
                        }, new UILoading.Loading { _time = 0, showType = 1 });
                    }

                    string selectedServerName = GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx].name;
                    Dictionary<string, object> properties = new Dictionary<string, object>();
                    properties["page"] = "登录界面";
                    properties["component"] = "创角界面";
                    properties["position"] = "";
                    properties["event_click"] = $"创建并登录(服务器：{selectedServerName})";
                    TDAnalyticsUtil.Track("click_event", properties, true);
                }
                else
                {
                    UITool.ShowTips(_creat.date);
                }
            }, 
            () =>{ UITool.ShowTips("失败"); });
    }

    private void OnRandmNameBtnClick()
    {
        int dicCount = CFG_Name.GetDic().Count;
        string nameStr = $"{CFG_Name.GetDic()[UnityEngine.Random.Range(1, dicCount)].Firstname}{CFG_Name.GetDic()[UnityEngine.Random.Range(1, dicCount)].Familyname}{CFG_Name.GetDic()[UnityEngine.Random.Range(1, dicCount)].Lastname}";
        GetView().input_Name.text = nameStr;

        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties["page"] = "登录界面";
        properties["component"] = "创角界面";
        properties["position"] = "";
        properties["event_click"] = "随机名字";
        TDAnalyticsUtil.Track("click_event", properties, true);
    }
    public override void OnEnter(params object[] args)
    {
        AudioManager.Instance.PlayBgMusic(AudioKey.Sound.BG_ChooseRole1);
        IsClick = false;
        OnJobChoiscClick(0);
        OnRandmNameBtnClick();
    }


    public override void OnExit()
    {
        AudioManager.Instance.PlayBgMusic(AudioKey.Sound.BG_Login);
    }
}