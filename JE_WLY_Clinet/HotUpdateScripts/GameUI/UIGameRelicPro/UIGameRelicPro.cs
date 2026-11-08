using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using HotUpdateScripts;
using JEngine.UI;
using System;

public class UIGameRelicPro : UIBase<UIModelGameRelicPro, UIViewGameRelicPro>
{
    public override UILayer GetLayer()
    {
        return UILayer.NoLayer;
    }
    public override void OnInit()
    {
        GetView().btn_Show.onClick.AddListener(delegate
        {
            OnUIShow(true);
        });

        GetView().btn_Close.onClick.AddListener(delegate
        {
            OnUIShow(false);
        });
    }

    public override void OnEnter(params object[] args)
    {
    }

    private void OnUIShow(bool isShow)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        GetView().obj_UI.SetActive(isShow);
        GetView().btn_Show.gameObject.SetActive(!isShow);
    }

    private Func<bool>[] ReliecCheak = {
        ()=>  GameRoomM.Instance.CheakMonstersDead(),
        ()=>  GameRoomM.Instance.CheakMyTeamHpHalf() ,
        ()=>  GameRoomM.Instance.CheakMyTeamAllLife(),
    };

    public void OnRelicInfoShow()
    {
        //if (GameRoomM.Instance._GameState == GameRoomM.GameState.Over)
        //    return;

        for (int i = 0; i < GetView().root_Tips.childCount; i++)
        {
            bool isUnlock = ReliecCheak[i].Invoke();
            Transform itemRoot = GetView().root_Tips.GetChild(i);
            JUIExtensions.txt(itemRoot).text = $"<color=#{(!isUnlock ? "FFFFFF" : "94E771")}>{GameConst.RelicItemStrs[i]}</color>";
            itemRoot.Find("on").gameObject.SetActive(isUnlock);
            itemRoot.Find("off").gameObject.SetActive(!isUnlock);
        }
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}