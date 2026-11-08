using System.Collections;
using UnityEngine;

public class BQSDKUILogicBase : MonoBehaviour
{
    public UIBQSDKMgr ParentCtr = null;
    public void Init(UIBQSDKMgr parent)
    {
        ParentCtr  = parent;
        AddEventListener();
        OnHide();
    }
    public virtual void OnShow() {
        this.gameObject.SetActive(true);
    }
    public virtual void OnHide() {
        this.gameObject.SetActive(false);
    }
    protected virtual void AddEventListener() { }
    protected virtual void RemoveEventListener() { }
    public void Exit() 
    {
        RemoveEventListener();
    }
}