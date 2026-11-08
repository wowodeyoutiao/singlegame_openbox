using MetaFramework.UI;
using Spine.Unity;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class UIViewNewPlayerGuide : MonoBehaviour,IUIView
{
	public GameObject maskGo;
    public Button skipBtn;

    public GoMainBattle goMainBattle;
    public GoPlayBattle goPlayBattle;
    public InFirstPlayBattleGuidance inFirstPlayBattleGuidance;
    public AddNewRoleInFormation addNewRoleInFormation;
    public CloseItemBoxAndGoBattleAgain closeItemBoxAndGoBattleAgain;

    [HideInInspector]
    public Dictionary<eGuideStep, NewPlayerGuideStepBase> guideSteps = new Dictionary<eGuideStep, NewPlayerGuideStepBase>();
    public void Awake()
    {
        
    }
	public void Init(GameObject go)
	{
        guideSteps.Clear();
        guideSteps.Add(eGuideStep.GoToMainBattle, goMainBattle);
        guideSteps.Add(eGuideStep.GoToBattle, goPlayBattle);
        guideSteps.Add(eGuideStep.InPlayBattleGuidance, inFirstPlayBattleGuidance);
        guideSteps.Add(eGuideStep.AddNewRoleInFormation, addNewRoleInFormation);
        guideSteps.Add(eGuideStep.ClocseItemBoxAndGoToMainFightAgain, closeItemBoxAndGoBattleAgain);

        foreach (var item in guideSteps)
        {
            if (item.Value == null)
                continue;
            item.Value.gameObject.SetActive(false);
        }
    }

    public void OnDestroy()
	{
		 UITool.RemoveUI<UINewPlayerGuide>();
	}

    ///步骤内执行标准
    //1、Steps[CurrentStep]?.SetActive(false);   先隐藏步骤内
    //2、UniGameEvent.AddEventListener<GameEvent.UIPanel_Operate>(WaitOpenUIFightChoice); //监听UI打开消息，   因为UI可能从缓存中直接打开。如果监听放前面或后面会引起奇怪的BUG
    //3、NewPlayerGuideWidget.Instance.InvokeButtonEvent("DownBtns_FightButton"); //点击按钮事件

}