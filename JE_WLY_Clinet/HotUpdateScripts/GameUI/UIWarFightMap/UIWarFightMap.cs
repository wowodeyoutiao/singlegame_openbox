using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.UI;
using Spine;
using HotUpdateScripts;

public class UIWarFightMap : UIBase<UIModelWarFightMap, UIViewWarFightMap>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIWarFightMap>();
        });

        GetView().btn_Help.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null,48);
        });
        for (int i = 0; i < CFG_Room.WarFights.Length; i++)
        {
            int uid = CFG_Room.WarFights[i].ID;
            GetView().Btns[i].btn.onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIWarFightInfo>(null, uid);
            });

        }
        OnWarFightUI();
    }

    public override void OnEnter(params object[] args)
    {
        TcpC.Get_C2S_WarRoomList();
    }

    public override void OnExit()
    {
        
    }

    private void OnWarFightUI()
    {
        for (int i = 0; i < CFG_Room.WarFights.Length; i++)
        {
            GetView().Btns[i].txt.text = CFG_Room.WarFights[i].Name;
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.WarFightAtlas, CFG_Room.WarFights[i].Name.Contains("困难")? "icon_qxzc_kn": "icon_qxzc_dy", GetView().Btns[i].img_type);
        }
    }

    public override void OnDestroy()
    {

    }
}