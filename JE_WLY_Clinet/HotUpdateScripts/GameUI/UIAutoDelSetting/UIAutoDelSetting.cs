using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;

public class UIAutoDelSetting : UIBase<UIModelAutoDelSetting, UIViewAutoDelSetting>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {

        GetView().dropdown_Qulite.AddOptions(GetModel().DroQualityList);
        GetView().dropdown_Qulite.onValueChanged.AddListener(SetSortQulite);

        GetView().dropdown_Level.AddOptions(GetModel().DroLevelList);
        GetView().dropdown_Level.onValueChanged.AddListener(SetSortlevel);

        GetView().btn_Power.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().obj_power_on.SetActive(!GetView().obj_power_on.activeSelf);
        });

        GetView().btn_Qulie.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().obj_Quile_on.SetActive(!GetView().obj_Quile_on.activeSelf);
        });

        GetView().btn_Lv.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().obj_Lv_on.SetActive(!GetView().obj_Lv_on.activeSelf);
        });

        GetView().btn_Close.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIAutoDelSetting>();
        });

        GetView().btn_Help.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 20);
        });

        GetView().btn_Sure.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            TcpC.C2S_AutoDisassembly(new pb.AutoDisassemblyInfo
            {
                Open = BagM.Instance.AutoDelSetInfo.Open,
                Level = select_Level,
                LevelSwitch = GetView().obj_Lv_on.activeSelf,
                Quality = (pb.ItemQuality)select_Qulite,
                QualitySwitch = GetView().obj_Quile_on.activeSelf,
                HighCombat = GetView().obj_power_on.activeSelf,
            }) ;

        });
    }

    private int select_Qulite;
    private void SetSortQulite(int idx)
    {
        select_Qulite = idx + 1;
    }

    private int select_Level;
    private void SetSortlevel(int idx)
    {
        select_Level = (idx + 1) * 10;
    }

    public override void OnEnter(params object[] args)
    {
        select_Level = BagM.Instance.AutoDelSetInfo.Level;
        select_Qulite = (int)BagM.Instance.AutoDelSetInfo.Quality;
        GetView().dropdown_Qulite.value = ((int)BagM.Instance.AutoDelSetInfo.Quality - 1);
        GetView().dropdown_Level.value = (BagM.Instance.AutoDelSetInfo.Level / 10) - 1;
        GetView().obj_Lv_on.SetActive(BagM.Instance.AutoDelSetInfo.LevelSwitch);
        GetView().obj_power_on.SetActive(BagM.Instance.AutoDelSetInfo.HighCombat);
        GetView().obj_Quile_on.SetActive(BagM.Instance.AutoDelSetInfo.QualitySwitch);
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}