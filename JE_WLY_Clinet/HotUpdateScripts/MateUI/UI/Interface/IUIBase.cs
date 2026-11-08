using UnityEngine;

namespace MetaFramework.UI
{
    public interface IUIBase
    {
        string uiName { get; set; }

        int _canvasOrder { get; set; }
        GameObject uiGo { get; set; }
        IUIModel uiModel { get; set; }
        IUIView uiView { get; set; }

        UILayer GetLayer();
        void OnInit();

        void OnLayerChange();

        void OnEnter(params object[] args);

        void OnExit();

        void OnDestroy();

    }
}
