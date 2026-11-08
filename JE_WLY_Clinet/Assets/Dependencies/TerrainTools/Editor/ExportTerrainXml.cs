using System;
using UnityEngine;
using System.Xml;
using System.IO;
using System.Text;  
using UnityEditor;
using UnityEngine.AI;
using UnityEngine.SceneManagement;
using System.Collections.Generic;
using System.Runtime.Serialization.Formatters.Binary;
using JEngine.LitJson;
using JEngine.Core;
using JEngine.Examples.Data;
using Palmmedia.ReportGenerator.Core.Common;

public class ExportTerrainXml   
{

    [System.Serializable]
    public class MapData
    {
        public int x;
        public int z;
        public decimal posY;

        public override string ToString()
        {
            return $"x={x},z={z},posY={posY}";
        }
    }

    public class mapJson<T>
    {
        public int mapWidth;
        public int mapHeight;
        public List<MapData> mapPosList;

        public override string ToString()
        {
            return $"mapWidth={mapWidth},mapHeight={mapHeight},mapPosList={mapPosList}";
        }
    }


    [MenuItem("TerrainTools/ExportTerrainXml")]
    public static void ExportXml()
    {
        Scene currentScene = SceneManager.GetActiveScene();
        // 打印当前场景的名称
        Debug.Log("Current Scene Name: " + currentScene.name);
        //if (currentScene.name.IndexOf("pve", StringComparison.Ordinal) <0 && currentScene.name.IndexOf("pvp", StringComparison.Ordinal) <0)
        //{
        //    Debug.LogError("当前场景不支持");
        //    return;
        //}

        // string tmpPath = Application.dataPath + "/" + currentScene.name + ".xml";
        // if (File.Exists(tmpPath))
        // {
        //     File.Delete(tmpPath);
        // }
        // StreamWriter tmpStreamWriter = new StreamWriter(tmpPath);

        //string str= exportPoint(100, 100,1);

        // tmpStreamWriter.Write(str);

        // tmpStreamWriter.Close();

        int widthAndHeight = 100;
        if (currentScene.name == "map_8")
            widthAndHeight = 30;
        else if (currentScene.name == "map_9")
            widthAndHeight = 50;
        else if (currentScene.name == "map_10")
            widthAndHeight = 20;


        exportPoint(widthAndHeight, widthAndHeight, currentScene.name, 0, 0);
    }


    public static void exportPoint(int width, int height, string currentMap, int offsetWid, int offsetHei)
    {
        StringBuilder str = new StringBuilder();

        List<string> listPos = new List<string>();
        List<MapData> mapDatas = new List<MapData>();

        //  listPos.Add(width + ":" + height + ":" + offsetWid + ":" + offsetHei);//宽高+偏移
        // str.Append("\r\n");
        for (int i = 0; i < width; ++i)
        {
            //str.Append("{");
            for (int j = 0; j < height; ++j)
            {
                int res = 1;
                NavMeshHit hit;

                float y = -100;

                // if (NavMesh.SamplePosition(new Vector3(i, 30, j), out hit, 50, NavMesh.AllAreas))
                // {
                //     res = 0;
                //     y = hit.position.y;
                // }else
                //     break;

                // if (i == 268)
                // {
                //
                //     //bool isnav  = NavMesh.SamplePosition(new Vector3(269, 7f, 306), out hit, 0.1f, NavMesh.AllAreas);
                //     
                //     Debug.Log("da");
                // }
                for (float k = 20; k >= -10; k -= 1.0f)
                {
                    // float posy = checkPossible(i, j, k);
                    // if (posy>=0)
                    // {
                    //     res = 0;
                    //     y = posy;
                    //     break;
                    // }

                    if (NavMesh.SamplePosition(new Vector3(i, k, j), out hit, 2.0f, NavMesh.AllAreas))
                    {
                        float hitx = hit.position.x;
                        float hity = hit.position.z;

                        if (Math.Floor(hitx) == i && Math.Floor(hity) == j)
                        {
                            float positionY = hit.position.y;
                            // if (positionY < 0)
                            //     positionY = 0.1f;
                            res = 0;
                            y = positionY;
                            break;
                        }
                    }
                }



                // Debug.DrawRay(startPos + new Vector3(j * accuracy, 0, -i * accuracy), Vector3.up, res == 1 ? Color.green : Color.red, 5);
               if (i==width-1&&j==height-1)
                    str.Append(res);
                
               else
                    str.Append(res).Append(",");

                if (y > -100)
                {
                    Debug.Log(i + ":" + y + ":" + j);
                    listPos.Add(i + ":" + y + ":" + j);

                    MapData _md = new MapData();
                    _md.x = i;
                    _md.z = j;
                    _md.posY = (decimal)y;

                    mapDatas.Add(_md);
                }
                // strClient.Append(i).Append(":").Append(y).Append(":").Append(j);

            }

          //  str.Append("\n");
        }

        Debug.Log(str.ToString());

        CreateXML(width, height, str.ToString(), currentMap);

        string path = Application.dataPath + "/" + currentMap + ".txt";
        BinaryFormatter bf = new BinaryFormatter();
       // string ResourcePath = System.IO.Path.Combine(tmpPath);
       // string path = ResourcePath + "/" + currentMap + ".txt";
        if (File.Exists(path))
            File.Delete(path);
        FileStream fileStream = System.IO.File.Create(path);


        bf.Serialize(fileStream, listPos);
        fileStream.Close();

        

        mapJson<MapData> _mapJson = new mapJson<MapData>()
        {
            mapWidth=width,
            mapHeight=height,
            mapPosList=mapDatas
            
        };

        var json = JsonMapper.ToJson(_mapJson);
        Log.Print($"序列化后的JSON: {json}");



        string pathJson = Application.dataPath + "/" + currentMap + ".json";
       // BinaryFormatter bfJson = new BinaryFormatter();
        
        if (File.Exists(pathJson))
            File.Delete(pathJson);
       // FileStream fileStreamJson = System.IO.File.Create(pathJson);
       // bf.Serialize(fileStreamJson, json);
       // fileStreamJson.Close();

       // string jsonStr = JsonUtility.ToJson(json);
        File.WriteAllText(pathJson, json);

        Log.Print("write json finish");
        // BinaryFormatter bfRead = new BinaryFormatter();
        // FileStream filestream = System.IO.File.Open(path, FileMode.Open);
        // List<string> mapataList = (List<string>)bf.Deserialize(filestream);
    }

    //public static void exportPoint(int wid, int hei)
    //{
    //    StringBuilder str = new StringBuilder();
    //    // StringBuilder strClient = new StringBuilder();
    //    List<string> listPos = new List<string>();
    //    listPos.Add(wid + ":" + hei + ":" + 33 + ":" + 520);//宽高+偏移
    //    //int[,] list = new int[wid, hei];
    //    // str.Append("startpos=").Append(startPos).Append("\r\n");
    //    //  str.Append("height=").Append(hei).Append("\r\nwide=").Append(wid).Append("\r\n");
    //    str.Append("\r\n");
    //    for (int i = 0; i < wid; ++i)
    //    {
    //        //str.Append("{");
    //        for (int j = 0; j < hei; ++j)
    //        {
    //            int res = 1;
    //            NavMeshHit hit;
    //            // if (i == 467 && j == 248)
    //            // {
    //            //     Debug.Log("stop");
    //            // }
    //            float y = 0;

    //            // if (i >= 396 && i <= 409 && j >= 420 && j <= 458)//刷成阻挡点
    //            // {
    //            //     
    //            // } 
    //            // else  if (i >= 375 && i <= 403 && j >= 749 && j <= 774)//刷成阻挡点
    //            // {
    //            //     
    //            // }
    //            // else
    //            // {

    //            for (float k = 25; k >= 15f; k -= 0.1f)
    //            {
    //                // NavMesh.
    //                // if (NavMesh.SamplePosition(new Vector3(i, k, j), out hit, 0.1f, NavMesh.AllAreas))
    //                // {
    //                //     res = 0;
    //                //     y = hit.position.y;
    //                //     if(i==51&&j==567)
    //                //     Debug.Log("y:"+y);
    //                //     break;
    //                // }

    //                float posy = checkPossible(i, j, k);
    //                if (posy >= 0)
    //                {
    //                    res = 0;
    //                    y = posy;
    //                    break;
    //                }
    //            }
    //            //  }



    //            // Debug.DrawRay(startPos + new Vector3(j * accuracy, 0, -i * accuracy), Vector3.up, res == 1 ? Color.green : Color.red, 5);
    //            str.Append(res).Append(",");

    //            if (y > 0)
    //                listPos.Add(i + ":" + y + ":" + j);
    //            // strClient.Append(i).Append(":").Append(y).Append(":").Append(j);
    //        }

    //        str.Append("\n");
    //    }

    //    Debug.Log(str.ToString());

    //    CreateXML(wid, hei, str.ToString(), 2);


    //    BinaryFormatter bf = new BinaryFormatter();
    //    string ResourcePath = System.IO.Path.Combine(PathUtil.GetWorkPath, "../mjDownData");
    //    string path = ResourcePath + "/map2.txt";
    //    if (File.Exists(path))
    //        File.Delete(path);
    //    FileStream fileStream = System.IO.File.Create(path);


    //    bf.Serialize(fileStream, listPos);
    //    fileStream.Close();


    //    // BinaryFormatter bfRead = new BinaryFormatter();
    //    // FileStream filestream = System.IO.File.Open(path, FileMode.Open);
    //    // List<string> mapataList = (List<string>)bf.Deserialize(filestream);
    //}

    public static float checkPossible(int wid, int hei, float y)
    {

        NavMeshHit hit = default;
        for (float i = wid; i < wid + 1; i += 0.1f)
        {
            for (float j = hei; j < hei + 1; j += 0.1f)
            {
                if (NavMesh.SamplePosition(new Vector3(i, y, j), out hit, 2f, NavMesh.AllAreas))
                {
                    float positionY = hit.position.y;
                    if (positionY < 0)
                        positionY = 0;
                    return positionY;//这个格子可行走

                }
                // else if (NavMesh.SamplePosition(new Vector3(i, y, j), out hit, 0.1f, NavMesh.AllAreas))
                // {
                //     return hit.position.y;//这个格子可行走
                // }
            }
        }

        return -1;
    }

    public static void CreateXML(int width, int height, string str, string currentMap)
    {
       // string ResourcePath = Application.dataPath + "/" + currentMap + ".xml";
        string path = Application.dataPath + "/" + currentMap + ".xml";
        if (File.Exists(path))
            File.Delete(path);

        if (!File.Exists(path))
        {
            //创建最上一层的节点。
            XmlDocument xml = new XmlDocument();
            XmlDeclaration xmldecl = xml.CreateXmlDeclaration("1.0", "UTF-8", "");
            xml.AppendChild(xmldecl);
            //创建最上一层的节点。
            XmlElement root = xml.CreateElement("map");
            root.SetAttribute("width", width.ToString());
            root.SetAttribute("height", height.ToString());
            root.SetAttribute("ID", currentMap.Replace("map_",""));
            //创建子节点
             XmlElement element = xml.CreateElement("values");
            //设置节点的属性 encoding="csv"
            // element.SetAttribute("encoding", "csv");
             element.InnerText = str;
             root.AppendChild(element);
            xml.AppendChild(root);
            //最后保存文件
            xml.Save(path);
        }
    }

}  