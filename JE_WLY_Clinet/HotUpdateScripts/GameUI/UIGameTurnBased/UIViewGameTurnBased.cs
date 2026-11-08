using HotUpdateScripts.Game;
using HotUpdateScripts.ProtoM;

using MetaFramework.UI;

using System.Collections;
using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

public class UIViewGameTurnBased : MonoBehaviour, IUIView
{
    public Transform Root;
    public Text txt_Time;
    public Button btn_Close;
    public Text txt_MonsterTip;
    // 我方标识
    public Transform root_MyTeam;
    // 对方标识
    public Transform root_OtherTeam;
    // 当前回合数和总回合数
    public Text txt_Rount;
    public Text txt_MyName;
    public Text txt_OtherName;

    public Button btn_JumpGameDisabled;
    public Button btn_JumpGame;
    public GameObject root_JumpTime;
    public Text txt_JumpTime;

    public void Awake()
    {

    }

    public void Init(GameObject go)
    {
        txt_Time = go.transform.Find("txt_Time").GetComponent<Text>();
        txt_MonsterTip = go.transform.Find("txt_Time/txt_MonsterTip").GetComponent<Text>();
    }

    public void OnDestroy()
    {
        UITool.RemoveUI<UIGameTurnBased>();
    }
}