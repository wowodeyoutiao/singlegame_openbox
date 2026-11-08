using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;

public class UIModelCreatRole : IUIModel
{
	public static Dictionary<JobType, string[]> JobDesDic = new Dictionary<JobType, string[]>
	{
		[JobType.Warrior] = new string[] { "冲锋陷阵", "坚不可摧", "血厚防高", "近战无敌" },
        [JobType.Mage] = new string[] { "毁天灭地", "雷霆万钧", "输出炸裂", "决胜千里" },
        [JobType.Sacrifice] = new string[] { "圣光普照", "众生庇佑", "群体增益", "妙手回春" },
    };
	
}