using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using MetaFramework.UI;

public class UIModelPartnerBook : IUIModel
{
    private List<string> _DroStrList;
    public List<string> DroStrList
    {
        get
        {
            if (_DroStrList == null)
            {
                _DroStrList = new List<string>(CFG_XMLConfig.xmlConfig.SelectQualityDic.Count);
                foreach (var item in CFG_XMLConfig.xmlConfig.SelectQualityDic)
                {
                    _DroStrList.Add(item.Value);
                }
            }
            return _DroStrList;
        }
    }
    private List<int> _DroQualityList;
    public List<int> DroQualityList
    {
        get
        {
            if (_DroQualityList == null)
            {
                _DroQualityList = new List<int>(CFG_XMLConfig.xmlConfig.SelectQualityDic.Count);
                foreach (var item in CFG_XMLConfig.xmlConfig.SelectQualityDic)
                {
                    _DroQualityList.Add(item.Key);
                }
            }
            return _DroQualityList;
        }
    }
}