using pb;
using System;
using UnityEngine;

namespace HotUpdateScripts.Game
{

    //[System.Serializable]
    //public class MapData
    //{
    //    public int x;
    //    public int z;
    //    public float posY;

    //}

    //[System.Serializable]
    //public class mapJson
    //{
    //    public int mapWidth;
    //    public int mapHeight;
    //    public List<MapData> mapPosList;
    //}

    public static class SceneManager
    {

        //private static Dictionary<int, Dictionary<Vector2, float>> mapPos; //地图可行走坐标
        //private static Dictionary<int, Vector2> mapWH; //地图宽高
        public static float Block = 1;

        public static int Wight = 0;
        public static int Height = 0;

        public static int MinX = 0;
        public static int MinY = 0;
        public static int MaxX = 0;
        public static int MaxY = 0;

        private static int OffX = 6;
        private static int OffY = 13;

        public static int EffOffZ = -3;
        public static void loadMapInfo(int mapID)
        {
            //if (mapPos == null)
            //    mapPos = new Dictionary<int, Dictionary<Vector2, float>>();
            //if (mapWH == null)
            //    mapWH = new Dictionary<int, Vector2>();

            //if (mapPos.ContainsKey(mapID))
            //    return;

            //_ = loadMap(mapID);

            Wight = CFG_Map.GetById(mapID).Wide;
            Height = CFG_Map.GetById(mapID).High;
            MinX = OffX;
            MaxX = Wight - OffX;

            MinY = OffY;
            MaxY = Height - OffY;
            if (MinY > MaxY)
            {
                MinY = Height / 2 - 1;
                MaxY = Height / 2 - 1;
            }
        }

        private static CFGItem_Map.BornItem[] _borns = null;

        //设置安全点
        public static void SetSafeBorns(int mapID)
        {
            _borns = mapID == 0 ? null : CFG_Map.GetById(mapID).Born;
        }

        /// <summary>
        /// 判断是否安全坐标
        /// </summary>
        /// <param name="x"></param>
        /// <param name="y"></param>
        /// <returns></returns>
        public static bool IsSafePos(double x, double y)
        {
            if (_borns != null)
            {
                for (int i = 0; i < _borns.Length; i++)
                {
                    if (getDistance2Pos(x, y, _borns[i].x, _borns[i].y) <= CFG_XMLConfig.xmlConfig.WarFight.bornRange)
                    {
                        return true;
                    }
                }
            }
            return false;
        }

        private static int getDistance2Pos(double x, double y, int x1, int y1)
        {
            return (int)Math.Max(Math.Abs(x - x1), Math.Abs(y - y1));
        }

        public static Vector3 MapPos()
        {
            return new Vector3(Wight / 2 + 1, Height / 2 + 1, 0);
        }

        //private static async Task loadMap(int mapFileName)
        //{
        //   // Dictionary<Point, float> mapPosList = new Dictionary<Point, float>();

        //    if (!mapPos.ContainsKey(mapFileName))
        //    {
        //        Dictionary<Vector2, float> _mapPosList = new Dictionary<Vector2, float>();
        //        mapPos.Add(mapFileName, _mapPosList);
        //    }

        //    await ResouseHelper.Instance.LoadTextAssetAsync("Other/TextAssets/mapPos/map_"+ mapFileName + ".json", (jsonStr) => {

        //        if (jsonStr == "")
        //        {
        //            Debug.Log("jerryc json is null");
        //        }
        //        mapJson _mapJson = JsonMapper.ToObject<mapJson>(jsonStr);

        //        int mapwid = _mapJson.mapWidth;
        //        Debug.Log("jerryc json wid:" + mapwid);

        //        int maphei = _mapJson.mapHeight;
        //        Debug.Log("jerryc json hei:" + maphei);

        //        if (!mapWH.ContainsKey(mapFileName))//地图宽 高
        //        {
        //            Vector2 _point = new Vector2();
        //            _point.x = mapwid;
        //            _point.y = maphei;
        //            mapWH.Add(mapFileName, _point);

        //        }

        //        foreach (var VARIABLE in _mapJson.mapPosList)
        //        {
        //            float x = VARIABLE.x;
        //            float z = VARIABLE.z;
        //            //float y = Convert.ToSingle(VARIABLE.posY);
        //            float y = VARIABLE.posY;
        //            Vector2 _pointPos = new Vector2();
        //            _pointPos.x = x;
        //            _pointPos.y = z;
        //            if (!mapPos[mapFileName].ContainsKey(_pointPos))
        //            {
        //                mapPos[mapFileName].Add(_pointPos, y);
        //            }
        //        }
        //    });

        //   // string resultString = $"{{\"jsonDatas\":{json.text}}}";
        //  //  Debug.Log("jerryc json" + json.text);


        //}


        public static float GetOffZ(double X, double Y)
        {
            if (X < 0 || X >= Wight || Y < 0 || Y >= Height)
                return Block;
            return (float)(200 - Y) * -0.01f;
        }

        //public static double getPosY(int mapID, double x, double z)
        //{
        //    Vector2 point = new Vector2();
        //    point.x = (float)x;
        //    point.y = (float)z;
        //    if (mapPos.ContainsKey(mapID))
        //    {
        //        //foreach(var pos in mapPos[mapID])
        //        //{
        //        //    Debug.Log("jerryc test pos key:" + pos.Key.x);
        //        //    Debug.Log("jerryc test pos key:" + pos.Key.y);
        //        //    Debug.Log("jerryc test pos  value:" + pos.Value);

        //        //    if (pos.Key.x == 50 && pos.Key.y == 47)
        //        //    {
        //        //        Debug.Log("aaaaa2222");
        //        //    }

        //        //}
        //        if (mapPos[mapID].ContainsKey(point))
        //        {
        //            return mapPos[mapID][point];
        //        }
        //    }
        //    else
        //        Debug.Log("mapID miss");

        //    return Block;
        //}


        public static Vector3 GetSeverRotationPos(pb.Direction dir)
        {
            Vector3 endPos = Vector3.zero;
            switch (dir)
            {
                case pb.Direction.UP:
                    endPos = new Vector3(0, 0, -90);
                    break;
                case pb.Direction.UP_RIGHT:
                    endPos = new Vector3(180, 0, 135);
                    break;
                case pb.Direction.RIGHT:
                    endPos = new Vector3(180, 0, 180);
                    break;
                case pb.Direction.DOWN_RIGHT:
                    endPos = new Vector3(180, 0, 225);
                    break;
                case pb.Direction.DOWN:
                    endPos = new Vector3(180, 0, 270);
                    break;
                case pb.Direction.DOWN_LEFT:
                    endPos = new Vector3(0, 0, 45);
                    break;
                case pb.Direction.LEFT:
                    endPos = Vector3.zero;
                    break;
                case pb.Direction.UP_LEFT:
                    endPos = new Vector3(0, 0, -45);
                    break;
                case pb.Direction.NONE:
                    break;
                default:
                    break;
            }
            return endPos;
        }

        public static Vector3 GetRotationPos(Direction dir,ref pb.Direction _Dir)
        {
            if (dir != Direction.UP && dir != Direction.DOWN)
                _Dir = dir;

            Vector3 endPos = Vector3.one;
            if (_Dir == Direction.UP_LEFT || _Dir == Direction.DOWN_LEFT || _Dir == Direction.LEFT)
                endPos = Vector3.zero;
            else if (_Dir == Direction.UP_RIGHT || _Dir == Direction.DOWN_RIGHT || _Dir == Direction.RIGHT)
                endPos = new Vector3(0, 180, 0);
            return endPos;
        }


    }
}
