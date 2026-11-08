using HotUpdateScripts;
using JEngine.Helper;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;


public class ComboItem : MonoBehaviour
{
    public Image imgCombo;
    public Image imgEffect;

    private Dictionary<int, string> comboNormalIconDic = new() { { 1, "btn_Main_fkan_2" }, { 2, "btn_Main_sjan_2" }, { 3, "btn_Main_yqan_2" }, { 4, "btn_Main_dcan_2" } };
    private Dictionary<int, string> comboHitIconDic = new() { { 1, "btn_Main_fkan_1" }, { 2, "btn_Main_sjan_1" }, { 3, "btn_Main_yqan_1" }, { 4, "btn_Main_dcan_1" } };

    public void Show(int comboNo)
    {
        this.gameObject.SetActive(true);

        if(comboNormalIconDic.ContainsKey(comboNo))
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.HallAtlas, comboNormalIconDic[comboNo], imgCombo);

        if (comboHitIconDic.ContainsKey(comboNo))
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.HallAtlas, comboHitIconDic[comboNo], imgEffect);

        imgEffect.gameObject?.SetActive(false);
    }

    public void Hit()
    {
        imgEffect.gameObject?.SetActive(true);
    }

    public void Hide()
    {
        this.gameObject.SetActive(false);
        imgEffect.gameObject?.SetActive(false);
    }
}