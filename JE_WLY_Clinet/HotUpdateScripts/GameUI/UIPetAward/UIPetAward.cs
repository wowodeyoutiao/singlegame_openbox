using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Helper;
using Spine.Unity;
using HotUpdateScripts;

public class UIPetAward : UIBase<UIModelPetAward, UIViewPetAward>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }


    public override void OnInit()
    {
        GetView().btn_back.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPetAward>();
        });
    }
    public override void OnLayerChange()
    {
        GameUtil.Instance.SetEffectLayer(GetView().obj_eff7, _canvasOrder );
        GameUtil.Instance.SetEffectLayer(GetView().obj_eff8, _canvasOrder);
    }
    public override void OnEnter(params object[] args)
    {
        var CFGPet = CFG_Pet.GetById((int)args[0]);
        GetView().obj_eff7.gameObject.SetActive(CFGPet.Quality == 7);
        GetView().obj_eff8.gameObject.SetActive(CFGPet.Quality == 8);
        GetView().txt_Name.text = GameUtil.Instance.GetStrByQuality(CFGPet.Quality, CFGPet.Name);

        if (GetView().root_Pet.childCount > 0)
        {
            PoolHelper.Instance.RecycObj(GetView().root_Pet.GetChild(0).gameObject);
        }

        string pooName = "Pet_" + CFGPet.ID;
        _ = GameHelper.Instance.SetAsyncObj(pooName, GameConst.GetMonsterModelByID(pooName), (_mainRole) =>
        {
            _mainRole.transform.SetParent(GetView().root_Pet);
            _mainRole.transform.localScale = Vector3.one * 1.5f;
            _mainRole.transform.localPosition = Vector3.zero;
            _mainRole.SetActive(true);
        });
    }

    public override void OnExit()
    {
        if (GetView().root_Pet.childCount > 0)
        {
            PoolHelper.Instance.RecycObj(GetView().root_Pet.GetChild(0).gameObject);
        }
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}