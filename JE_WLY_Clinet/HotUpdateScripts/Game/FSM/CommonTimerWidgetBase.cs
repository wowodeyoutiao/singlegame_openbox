using System.Collections;
using System.Collections.Generic;
using UnityEngine;

/// <summary>
///  通用定时器基类
/// </summary>
/// <typeparam name="TDevice"></typeparam>
public class CommonTimerWidgetBase<TDevice> : MonoBehaviour where TDevice : CommonTimerWidgetBase<TDevice>
{
    public static GameObject Root { get; private set; } = null;
    public static Transform RootTran { get; private set; } = null;
    public static TDevice Self { get; private set; } = null;
    public static string ActuatorTag { get; private set; } = null;
    protected static void Start(string _ActuatorTag, Transform _Dock = null)
    {
        if (Root != null) return;
        if (string.IsNullOrEmpty(_ActuatorTag)) _ActuatorTag = typeof(TDevice).Name;
        ActuatorTag = _ActuatorTag;
        GameObject _Go = new GameObject(_ActuatorTag);
        Root = _Go;
        RootTran = _Go.transform;

        if (_Dock != null)
        {
            RootTran.SetParent(_Dock);
            RootTran.localPosition = Vector3.zero;
            RootTran.localEulerAngles = Vector3.zero;
        }

        Self = _Go.AddComponent<TDevice>();
        Self.enabled = true;
    }

    protected static void Stop()
    {
        if (Root == null) return;
        Self.enabled = false;
        GameObject.Destroy(Root);
        Root = null;
        RootTran = null;
        Self = null;
    }
}
