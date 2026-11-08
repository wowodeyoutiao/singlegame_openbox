using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Helper;
using JEngine.UI;
using HotUpdateScripts.ProtoM;
using HotUpdateScripts.Game;

public class UIPVPGameOver : UIBase<UIModelPVPGameOver, UIViewPVPGameOver>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_yellow.onClick.AddListener(()=> {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            ClearJa();
            OnBack();
        });
    }

    private void OnBack()
    {
        if (GameRoomM.Instance._GameState == GameRoomM.GameState.No)
            UITool.Close<UIPVPGameOver>();
        else
            FightController.Instance.OnExitGame();
    }
    public override void OnEnter(params object[] args)
    {
        bool isWin = (bool)args[0];
        GetView().obj_Lose.SetActive(!isWin);
        GetView().obj_Win.SetActive(isWin);
        var Awrads = isWin ? CFG_XMLConfig.xmlConfig.PVPInfo.winAwards : CFG_XMLConfig.xmlConfig.PVPInfo.loseAwards;

        for (int i = 0; i < Awrads.Length; i++)
        {
            var go = GetView().root_Awards.GetChild(i);
            var cfg_item = Awrads[i];
            GameUtil.Instance.SetItemIcon(cfg_item.id, JUIExtensions.img(go.transform.Find("icon")), JUIExtensions.img(go.transform));

            int _num = cfg_item.mark == 0 ? cfg_item.num : cfg_item.num * RoleM.Instance.Level;
            JUIExtensions.txt(go.transform.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(_num);

            JUIExtensions.btn(go).onClick.RemoveAllListeners();
            JUIExtensions.btn(go).onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfg_item.id });
            });
        }

        GetView().txt_winJF.text = "积分+" + GameUtil.Instance.FormatNumWithUnit((int)args[1]);
        GetView().txt_loseJF.text = "积分" + GameUtil.Instance.FormatNumWithUnit((int)args[1]);

        int Count = 5;
        Ja = GameHelper.Instance.AddInvoke(new TimeItem(Count, 1, () => {

            GetView().txt_Time.text = Count + "秒后将自动退出";
            Count--;
            if (Count<=0)
                OnBack();
        }));
       
    }

    private TimeItem Ja;
    private void ClearJa()
    {
        if (Ja != null && GameHelper.Instance != null)
        {
            GameHelper.Instance.RemoveInvoke(Ja);
            Ja = null;
        }
    }

    public override void OnExit()
    {
        ClearJa();
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}