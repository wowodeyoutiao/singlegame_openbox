using UnityEngine;
using UnityEngine.UI;

public class CanvasScale : MonoBehaviour
{

    private void Start()
    {
        var canvasScaler = this.transform.GetComponent<CanvasScaler>();
        var resolution = canvasScaler.referenceResolution;
        var screenRatio = Screen.width / Screen.height;
        if (screenRatio > (resolution.x / resolution.y))
            canvasScaler.matchWidthOrHeight = 1;
        else
            canvasScaler.matchWidthOrHeight = 0;
    }


}
