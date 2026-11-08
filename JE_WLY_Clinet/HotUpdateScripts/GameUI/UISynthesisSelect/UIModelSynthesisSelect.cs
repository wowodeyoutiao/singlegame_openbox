using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using MetaFramework.UI;

public class UIModelSynthesisSelect : IUIModel
{
	public List<string> DroSexList = new List<string> { CFG_Item.GetSexStrById(0), CFG_Item.GetSexStrById(1), CFG_Item.GetSexStrById(2) };
	public List<string> DroJobList = new List<string> { CFG_Item.GetJobStrById(0), CFG_Item.GetJobStrById(1), CFG_Item.GetJobStrById(2), CFG_Item.GetJobStrById(3) };
	public List<string> DroLevelList 
	{
		get {
			List<string> list = new List<string>(CFG_XMLConfig.xmlConfig.Sift_Level.Length);
			foreach (var item in CFG_XMLConfig.xmlConfig.Sift_Level)
			{
				list.Add(item.Option);
			}
			return list;
		}
	}

	public List<string> DroQualityList
	{
		get
		{
			List<string> list = new List<string>(CFG_XMLConfig.xmlConfig.Sift_Quality.Length);
			foreach (var item in CFG_XMLConfig.xmlConfig.Sift_Quality)
			{
				list.Add(item.Option);
			}
			return list;
		}
	}

}