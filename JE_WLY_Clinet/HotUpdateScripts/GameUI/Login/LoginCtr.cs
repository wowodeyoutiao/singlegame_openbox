using HotUpdateScripts.ProtoM;
using MetaFramework.UI;
using System;
using UnityEngine;

namespace HotUpdateScripts
{
    public class LoginCtr : MonoBehaviour
    {
        public void Awake()
        {
            if (GameM.Instance.SeverTime == 0)
                GameM.Instance.SeverTime = (double)((DateTime.Now.ToUniversalTime() - new DateTime(1970, 1, 1, 0, 0, 0, DateTimeKind.Utc)).TotalSeconds);
            UITool.Open<UILogin>();
        }
    }
}


