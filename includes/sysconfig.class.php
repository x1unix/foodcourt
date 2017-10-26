<?php
	// Developed by Stepan Kupyak <skupyak@lohika.com>
	// Edited by kosyak <kosyak_ua@yahoo.com>
	
	require_once (DB_FOR_USE."_core.class.php");
	require_once("user.class.php");
	require_once("validator.class.php");

	class SysConfig{
		var $data;
		var $error = array();
		var $db;

		function SysConfig(){
			$className = DB_FOR_USE.'Core';
			$this->db = new $className;
		}

		function saveConfig($config){
			$data['table_name'] = SYS_CONFIG_TBL;
			$data['data'] = $config;

			if ($this->db->add_data($data)){
				return true;
			}else{
				$this->error = $this->db->error;
				return false;
			}
		}

		function editConfig($id, $config){
			$table = SYS_CONFIG_TBL;
			$status = $this->db->edit_data(array(
				'where_clause'=> "id=$id",
				'data'=>$config,
				'table_name'=> $table,
			));
			if (!$status){
				$this->error = $this->db->error;
			}
			return $status;
		}

		function loadConfig($id = false){
			$this->data = array();
			$table = SYS_CONFIG_TBL;
			$id = mysql_real_escape_string($id);
			if ($id){
				$where = "id = $id";
			}
			$status = $this->db->get_data(array(
				'where_clause'=> $where,
				'table_name'=> $table,
			));
			if (!$status){
				$this->error = $this->db->error;
				return false;
			}

			foreach ($this->db->data as $value){
				$this->data[$value['id']] = $value;
			}
			return true;
		}

		function deleteConfig($id){
			$table = SYS_CONFIG_TBL;
			$id = mysql_real_escape_string($id);
			if ($id){
				$where = "id=$id";
			}
			$status = $this->db->delete_data(array(
				'where_clause'=> $where,
				'table_name'=> $table,
			));
			if (!$status){
				$this->error = $this->db->error;
			}
			$this->data = $this->db->data;
			return true;
		}

		function getCategories(){
			$table = SYS_CATEGORIES_TBL;
			$status = $this->db->get_data(array(
				'table_name'=> $table,
			));
			if (!$status){
				$this->error = $this->db->error;
				return false;
			}
			$data = $this->db->data;

			$arrayData = array();

			foreach ($data as $value){
				$arrayData[$value['id']] = $value;
			}
			$this->data = $arrayData;
			return true;
		}

		function getConfigVariables($all = true){
			$this->loadConfig();
			$arrayData = array();

			foreach ($this->data as $value){
				$arrayData[$value['id']] = $value;
			}

			return $arrayData;
		}

		function addCategory($name, $localization_var){
			$table = SYS_CATEGORIES_TBL;
			$name = mysql_real_escape_string($name);
			$localization_var = mysql_real_escape_string($localization_var);

			$data['table_name'] = $table;
			$data['data'] = array('name'=>$name, 'localization_var'=>$localization_var);

			if ($this->db->add_data($data)){
				return true;
			}else{
				$this->error = $this->db->error;
				return false;
			}
		}

		function editCategory($id, $name, $localization_var){
			$table = SYS_CATEGORIES_TBL;

			$status = $this->db->edit_data(array(
				'where_clause'=> "id=$id",
				'data'=> array('name'=>$name, 'localization_var'=>$localization_var),
				'table_name'=> $table,
			));
			if (!$status){
				$this->error = $this->db->error;
			}
			return $status;
		}

		function deleteCategory($id){
			$id = mysql_real_escape_string($id);
			$table = SYS_CATEGORIES_TBL;
			$where = "`id`='$id'";

			$status = $this->db->delete_data(array(
				'where_clause'=> $where,
				'table_name'=> $table,
			));

			if (!$status){
				$this->error[] = mysql_error();
				return false;
			}else{
				$table = SYS_CONFIG_TBL;

				$status = $this->db->edit_data(array(
					'where_clause'=> "cat_selected=$id",
					'data'=> array('cat_selected'=>"0"),
					'table_name'=> $table,
				));
				if (!$status){
					$this->error[] = mysql_error();
					return false;
				}
			}
			return true;
		}

		function updateVariables($valueArray){
			foreach ($valueArray as $id=>$array){
				$id = mysql_real_escape_string($id);
				$status = $this->editConfig($id, $array);
				if (!$status){
					$this->error[] = mysql_error();
					return false;
				}
			}
			return true;
		}

		function assignVariables()
		{
			# Stub
			$config = $this->getConfigVariables(false);

			foreach ($config as $value)
			{
				if ($value['active'])
				{
					define($value['name'], $value['value']);
				}
			}
			return true;
		}
	}
?>
