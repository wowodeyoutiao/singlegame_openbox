using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.Game;
using JEngine.Core;
using pb;
using JEngine.UI;
using JEngine.Helper;

public class UIGVGGameOver : UIBase<UIModelGVGGameOver, UIViewGVGGameOver>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_back.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            ClearJa();
            OnBack();
        });
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

    private void OnBack()
    {
        FightController.Instance.OnExitGame();
    }

    public override void OnEnter(params object[] args)
    {
        var info = (List<MeleeResult>)args[0];

        for (int i = 0; i < info.Count; i++)
        {
            Transform root = GetView().root_Rank.GetChild(i);
            JUIExtensions.txt(root.Find("txt_Name")).text = info[i].Player.Nickname;
            JUIExtensions.txt(root.Find("txt_killNum")).text = info[i].KillNum + "人";
            root.Find("Awards").gameObject.SetActive(info[i].Items.Count > 0);
            root.Find("txt_No").gameObject.SetActive(info[i].Items.Count == 0);

            if (info[i].Items.Count > 0)
            {
                for (int j = 0; j < 3; j++)
                {
                    root.Find("Awards").GetChild(j).gameObject.SetActive(info[i].Items.Count > j);
                    if (info[i].Items.Count > j)
                    {
                        JUIExtensions.txt(root.Find("Awards").GetChild(j).Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(info[i].Items[j].Num);
                        GameUtil.Instance.SetItemIcon(info[i].Items[j].CfgId, JUIExtensions.img(root.Find("Awards").GetChild(j).Find("icon")), JUIExtensions.img(root.Find("Awards").GetChild(j)));
                    }
                }
            }

        }

        int Count = 5;
        Ja = GameHelper.Instance.AddInvoke(new TimeItem(Count, 1, () => {
            if (isNullGo) return;
            GetView().txt_Time.text = Count + "秒后将自动退出";
            Count--;
            if (Count == 0)
                OnBack();
        }));
            
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