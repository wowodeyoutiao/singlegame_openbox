using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using JEngine.Helper;
using HotUpdateScripts.ProtoM;

public class UIPrivilegeCard : UIBase<UIModelPrivilegeCard, UIViewPrivilegeCard>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPrivilegeCard>();
        });
    }

    Color bule => new Color(0.2f, 0.29f, 0.58f);
    Color yellow => new Color(0.44f, 0.25f, 0.1f);
    Color purple => new Color(0.33f, 0.15f, 0.52f);

    private void OnUIInit()
    {
        for (int i = 0; i < GetView().root_Base.childCount; i++)
        {
            var item = GetView().root_Base.GetChild(i);
            var cfgInfo = CFG_Privilege.GetById(i + 1);
            JUIExtensions.txt(item.Find("desc/txt")).text = cfgInfo.Pdisplay;

            if(i==0)
                JUIExtensions.txt(item.Find("desc/txt")).color = bule;
            else if(i==1)
                JUIExtensions.txt(item.Find("desc/txt")).color = yellow;
            else if(i==2)
                JUIExtensions.txt(item.Find("desc/txt")).color = purple;

            int cfg_ID = cfgInfo.Dailyaward[0].id;

            JUIExtensions.btn(item.Find("img_award/Award")).onClick.RemoveAllListeners();
            JUIExtensions.btn(item.Find("img_award/Award")).onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfg_ID });
            });

            GameUtil.Instance.SetItemIcon(cfg_ID, JUIExtensions.img(item.Find("img_award/Award/icon")));
            JUIExtensions.txt(item.Find("img_award/Award/txt_num")).text = GameUtil.Instance.FormatNumWithUnit(cfgInfo.Dailyaward[0].num);
            GameUtil.Instance.SetItemIcon(CFG_Store.GetById(cfgInfo.SID).Solditems[0].itemID, JUIExtensions.img(item.Find("img_award/txt_coin/img_coin")));
            JUIExtensions.txt(item.Find("img_award/txt_coin")).text = GameUtil.Instance.FormatNumWithUnit(CFG_Store.GetById(cfgInfo.SID).Solditems[0].Num);
            JUIExtensions.txt(item.Find("btn_Buy/txt")).text = CFG_Store.GetById(cfgInfo.SID).Limite.Num + "元";
            JUIExtensions.txt(item.Find("img_time/txt_Time")).text = $"持续{cfgInfo.Duration}天";
            
            JUIExtensions.btn(item.Find("btn_Buy")).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

                System.Action JumpFun = delegate { UITool.Close<UIPrivilegeCard>(); };
                if (!GameUtil.Instance.ComCoinToCheck(cfgInfo.SID, JumpFun))
                {
                    return;
                }
                GameUtil.Instance.GameToBuy(cfgInfo.SID,()=> {
                    OnHttpRefresh();
                });
            });

            JUIExtensions.btn(item.Find("btn_Get")).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                HttpC.GetPrivilegeSign(cfgInfo.PID, (_r) =>{
                    OnHttpRefresh();
                } );
            });
        }
    }

    private void OnHttpRefresh()
    {
        HttpC.GetPrivilegeList((_r) => {
            if (!isNullGo)
                OnUIRefrsh(_r.data);
        });
    }

    private void OnUIRefrsh(HTTP_PrivilegeDetailDay[] data)
    {
        Dictionary<int, HTTP_PrivilegeDetailDay> pDic = new Dictionary<int, HTTP_PrivilegeDetailDay>();
        for (int i = 0; i < data.Length; i++)
        {
            pDic.Add(data[i].pid, data[i]);
        }

        for (int i = 0; i < GetView().root_Base.childCount; i++)
        {
            var item = GetView().root_Base.GetChild(i);
            var cfgInfo = CFG_Privilege.GetById(i + 1);

            item.Find("btn_Buy").gameObject.SetActive(pDic[cfgInfo.PID].buy == 1);
            item.Find("btn_Get").gameObject.SetActive(pDic[cfgInfo.PID].buy == 2 && pDic[cfgInfo.PID].get == 1);
            item.Find("img_award/Award/geted").gameObject.SetActive(pDic[cfgInfo.PID].buy == 2 && pDic[cfgInfo.PID].get == 2);
            item.Find("txt_Geted").gameObject.SetActive(pDic[cfgInfo.PID].buy == 2 && pDic[cfgInfo.PID].get == 2);
            if (pDic[cfgInfo.PID].buy != 1)
            {
                int day = GameUtil.Instance.GetTimeToD(pDic[cfgInfo.PID].end_time - GameUtil.Instance.GetCurrentTimeSeconds());
                string timeStr = day < 1 ? "不足1天" : day + "天";
                JUIExtensions.txt(item.Find("img_time/txt_Time")).text = $"{timeStr}";
            }
        }
    }

    public override void OnEnter(params object[] args)
    {
        RefreshUI();
    }

    private void RefreshUI()
    {
        OnUIInit();
        OnHttpRefresh();
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}