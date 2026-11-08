using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using MetaFramework.UI;

public class UIModelAutoDelSetting : IUIModel
{
	public List<string> DroLevelList
	{
		get
		{
			List<string> list = new List<string>(CFG_XMLConfig.xmlConfig.Sift_Level.Length);
            for (int i = 1; i < CFG_XMLConfig.xmlConfig.Sift_Level.Length; i++)
            {
				list.Add(CFG_XMLConfig.xmlConfig.Sift_Level[i].Option);

			}
			return list;
		}
	}

	public List<string> DroQualityList
	{
		get
		{
			List<string> list = new List<string>(CFG_XMLConfig.xmlConfig.Sift_Quality.Length);
			for (int i = 1; i < CFG_XMLConfig.xmlConfig.Sift_Quality.Length; i++)
			{
				list.Add(CFG_XMLConfig.xmlConfig.Sift_Quality[i].Option);

			}
			return list;
		}
	}
}