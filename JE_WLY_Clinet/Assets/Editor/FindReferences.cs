using UnityEngine;
using System.Collections;
using UnityEditor;
using System.IO;
using System.Linq;
using System.Text.RegularExpressions;
using System.Collections.Generic;

/// <summary>
/// Find All References on Unity Project.
/// </summary>
public class FindReferences
{
    enum ResourceSuffixType
    {
        All, Prefab, Unity, Mat, Asset
    }

    [MenuItem("Assets/Find References/*.mat", false, 11)]
    static private void FindInMats()
    {
        Find(ResourceSuffixType.Mat);
    }
    [MenuItem("Assets/Find References/*.prefab", false, 12)]
    static private void FindInPrefabs()
    {
        Find(ResourceSuffixType.Prefab);
    }
    [MenuItem("Assets/Find References/*.unity", false, 13)]
    static private void FindInUnitys()
    {
        Find(ResourceSuffixType.Unity);
    }
    [MenuItem("Assets/Find References/*.asset", false, 14)]
    static private void FindInAssets()
    {
        Find(ResourceSuffixType.Asset);
    }
    [MenuItem("Assets/Find References/*.*", false, 15)]
    static private void FindInAll()
    {
        Find(ResourceSuffixType.All);
    }
    static private void Find(ResourceSuffixType type)
    {
        EditorSettings.serializationMode = SerializationMode.ForceText;
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        if (!string.IsNullOrEmpty(path))
        {
            List<string> withExtensions;
            string guid = AssetDatabase.AssetPathToGUID(path);
            switch (type)
            {
                case ResourceSuffixType.Prefab:
                    withExtensions = new List<string>() { ".prefab" };
                    break;
                case ResourceSuffixType.Mat:
                    withExtensions = new List<string>() { ".mat" };
                    break;
                case ResourceSuffixType.Unity:
                    withExtensions = new List<string>() { ".unity" };
                    break;
                case ResourceSuffixType.Asset:
                    withExtensions = new List<string>() { ".asset" };
                    break;
                default:
                    withExtensions = new List<string>() { ".prefab", ".unity", ".mat", ".asset" };
                    break;
            }
            string[] files = Directory.GetFiles(Application.dataPath, "*.*", SearchOption.AllDirectories)
                .Where(s => withExtensions.Contains(Path.GetExtension(s).ToLower())).ToArray();
            int startIndex = 0;

            EditorApplication.update = delegate ()
            {
                string file = files[startIndex];
                bool isCancel = EditorUtility.DisplayCancelableProgressBar("匹配资源中", file, (float)startIndex / (float)files.Length);
                if (Regex.IsMatch(File.ReadAllText(file), guid))
                {
                    Debug.Log(file, AssetDatabase.LoadAssetAtPath(GetRelativeAssetsPath(file), typeof(object)));
                }

                startIndex++;
                if (isCancel || startIndex >= files.Length)
                {
                    EditorUtility.ClearProgressBar();
                    EditorApplication.update = null;
                    startIndex = 0;
                    Debug.Log("匹配结束");
                }
            };
        }
    }

    [MenuItem("Assets/Find References", true)]
    static private bool VFind()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        return (!string.IsNullOrEmpty(path));
    }

    static private string GetRelativeAssetsPath(string path)
    {
        return "Assets" + Path.GetFullPath(path).Replace(Path.GetFullPath(Application.dataPath), "").Replace('\\', '/');
    }
}