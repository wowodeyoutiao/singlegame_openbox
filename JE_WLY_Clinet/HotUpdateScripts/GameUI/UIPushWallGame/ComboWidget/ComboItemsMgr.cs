using JEngine.Core;
using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;


public class ComboItemsMgr : MonoBehaviour
{
    public Transform tfComboItemsRoot;
    public Text txtTips;

    public Sprite spCombo_1;
    public Sprite spCombo_2;
    public Sprite spCombo_3;
    public Sprite spCombo_4;

    //存储comboNo和对应的外观信息
    //private Dictionary<int,Sprite> spriteCollect = new(); 

    //存储UI上的每个combo信息
    private List<ComboItem> comboItems = new();

    //匹配结果回调
    private Action<bool> matchComboResult = null;

    private List<int> comboSequence = new(); //combo队列信息
    private int curentNo = 0;

    public void Awake()
    {
        //addSprite(1, spCombo_1);
        //addSprite(2, spCombo_2);
        //addSprite(3, spCombo_3);
        //addSprite(4, spCombo_4);

        comboItems.Clear();

        for (int i = 0; i < tfComboItemsRoot.childCount; i++)
        {
            ComboItem item = tfComboItemsRoot.GetChild(i).GetComponent<ComboItem>();

            if (item == null)
            {
                tfComboItemsRoot.GetChild(i).gameObject.SetActive(false);
                Debug.LogError($"在{tfComboItemsRoot.name}下的第{i+1}个的ComboItem组件不能为空！");
                continue;
            }

            item.Hide();
            comboItems.Add(item);
        }
    }

    //private void addSprite(int comboNo,Sprite spTexture)
    //{
    //    if (spTexture == null)
    //    {
    //        Debug.LogError($"第{comboNo}个的Sprite不能为空！");
    //    }

    //    spriteCollect[comboNo] = spTexture;
    //}

    /// <summary>
    /// 初始化
    /// </summary>
    /// <param name="m_comboResult"></param>
    public void Init(Action<bool> m_comboResult)
    {
        matchComboResult = m_comboResult;
    }

    public void ShowCombo(int count)
    {
        gameObject?.SetActive(true);

        if (count >= comboItems.Count)
            count = comboItems.Count;

        //初始化
        curentNo = 0;
        comboSequence.Clear();

        for (int i = 0; i < count; i++)
        {
            int comboNo = UnityEngine.Random.Range(1, 5); //随出Combo的值（取1，2，3，4这4个值）
            comboSequence.Add(comboNo);  //放入队列中

            //comboItems[i].Show(spriteCollect[comboNo]); //在UI上显示Combo的样子
            comboItems[i].Show(comboNo); //在UI上显示Combo的样子
        }
    }

    //匹配
    public void Match(int comboNo)
    {
        if (comboSequence.Count == 0) return;

        if (comboSequence[curentNo] != comboNo)//单个匹配失败
        {
            comboSequence.Clear();
            matchComboResult?.Invoke(false);
            return;
        }

        curentNo ++;

        if (curentNo >= comboSequence.Count) //完全匹配
        {
            comboSequence.Clear();
            matchComboResult?.Invoke(true);
            return;
        }
    }

    //是否接受输入
    public bool IsCanInput()
    {
        return comboSequence.Count > 0;
    }

    public void Hide()
    {
        comboSequence.Clear();
        curentNo = 0;

        gameObject?.SetActive(false);
    }
}
