using JEngine.Core;
using System.Collections;
using System.Text;
using UnityEngine;
using UnityEngine.Networking;

public class NetworkHelper : MonoBehaviour
{
    public enum HTTPMethod 
    {
        POST,
        GET,
        DELETE,
        PUT,
    }
    // Start is called before the first frame update
    public static void SendPost(string url, HTTPMethod htpMethod, string postStr, System.Action<string> onSuccess, System.Action<long, string> onFailure)
    {
        Debug.Log("Post " + postStr);
        CoroutineMgr.Instance.StartCoroutine(PostRequest(url, htpMethod, postStr, onSuccess, onFailure));
    }

    private static IEnumerator PostRequest(string url, HTTPMethod htpMethod, string postStr, System.Action<string> onSuccess, System.Action<long,string> onFailure) 
    {
        UnityWebRequest request = new UnityWebRequest(url, htpMethod.ToString());
        byte[] postData = Encoding.UTF8.GetBytes(postStr);
        request.uploadHandler = new UploadHandlerRaw(postData);
        request.downloadHandler = new DownloadHandlerBuffer();
        request.SetRequestHeader("Content-Type", "application/json");

        yield return request.SendWebRequest();

        if (request.result != UnityWebRequest.Result.Success || request.responseCode != 200)
        {
            Debug.LogError("Request Error: " + request.error);
            onFailure?.Invoke(request.responseCode, request.downloadHandler.text);
        }
        else
        {
            onSuccess?.Invoke(request.downloadHandler.text);
        }

        Debug.Log("httpResponse.StatusCode: " + request.responseCode);
    }

    public static void SendPost(string url, HTTPMethod htpMethod,  string headStr,string postStr, System.Action<string> onSuccess, System.Action<long, string> onFailure)
    {
       //iz Debug.Log(htpMethod +" --" + postStr);
        CoroutineMgr.Instance.StartCoroutine(PostRequestCoroutine(url, htpMethod, headStr, postStr, onSuccess, onFailure));
    }

    private static IEnumerator PostRequestCoroutine(string url, HTTPMethod htpMethod,  string headStr, string postStr, System.Action<string> onSuccess, System.Action<long, string> onFailure)
    {
        UnityWebRequest request = new UnityWebRequest(url, htpMethod.ToString());

        request.SetRequestHeader("Authorization", $"Bearer {headStr}");

        if (!string.IsNullOrEmpty(postStr))
        {
            request.SetRequestHeader("Content-Type", "application/json");
            byte[] postData = Encoding.UTF8.GetBytes(postStr);
            request.uploadHandler = new UploadHandlerRaw(postData);
        }

        request.downloadHandler = new DownloadHandlerBuffer();
        yield return request.SendWebRequest();

        if (request.result == UnityWebRequest.Result.Success && request.responseCode == 200)
        {
            onSuccess?.Invoke(request.downloadHandler.text);
        }
        else
        {
            // Debugging
            Debug.LogWarning($"Request failed. Error: {request.error}, Response Code: {request.responseCode}");
            onFailure?.Invoke(request.responseCode, request.downloadHandler.text);
        }
    }
}
