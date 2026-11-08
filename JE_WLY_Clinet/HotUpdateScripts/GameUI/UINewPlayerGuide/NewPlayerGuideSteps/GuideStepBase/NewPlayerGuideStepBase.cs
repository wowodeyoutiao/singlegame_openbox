using DG.Tweening;
using HotUpdateScripts.ProtoM;
using JEngine.Core;
using MetaFramework.UI;
using pb;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using UnityEditor;
using UnityEngine;
using UnityEngine.UI;


public abstract class NewPlayerGuideStepBase : MonoBehaviour
{
    [HideInInspector]
    protected List<GameObject> Steps = new List<GameObject>();
    [HideInInspector]
    protected int CurrentStep = -1;
    [HideInInspector]
    protected UIBase<UIModelNewPlayerGuide, UIViewNewPlayerGuide> ParentController;

    protected object[] Data { get; set; }

    //初始化
    public virtual void Enter(UIBase<UIModelNewPlayerGuide, UIViewNewPlayerGuide> parentController)
    { 
        this.ParentController = parentController;

        CurrentStep = -1;
        Steps.Clear();
        for (int i = 0; i < transform.childCount; i++)
        {
            Steps.Add(transform.GetChild(i).gameObject);
            Steps[Steps.Count - 1].SetActive(false);
        }
    }

    public virtual void DoNextStep()
    {
        if (CurrentStep > -1 && CurrentStep < Steps.Count)
            Steps[CurrentStep].SetActive(false);

        CurrentStep ++;

        if (CurrentStep >= Steps.Count)
        {
            Exit();
            return;
        }

        Steps[CurrentStep].SetActive(true);
    }

    /// <summary>
    /// 用于断线重连
    /// </summary>
    public abstract void ReEnter();

    public abstract void Exit();

    public void FlyButton(Transform tfSelfTrans,Vector3 beginPos,string targetPosName,float time = 1.0f, Action callBack = null)
    {
        tfSelfTrans.position = beginPos;
        (bool isOk, Vector3 wordPos) = NewPlayerGuideWidget.Instance.GetPosition(targetPosName);
        if (isOk)
        {
            tfSelfTrans.DOMove(wordPos, time).OnComplete(() => { callBack?.Invoke(); });
        }
        else
        {
            callBack?.Invoke();
        }
    }
}