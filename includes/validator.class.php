<?php
	require_once (DB_FOR_USE."_core.class.php");

	class Validator{
		function Validator(){
			$className = DB_FOR_USE.'Core';
			$this->db = new $className;
		}
		function validateType($var, $type){
			if ($type=="bool"){
				return in_array(strtolower($var), array('false', 'true'));
			}
			elseif ($type=="int"){
				return is_numeric($var);
			}
			else{
				return true;
			}
		}

		function validateName($name, $id, $table){
			$data = array();
			$name = mysql_real_escape_string($name);
			$where = "name='$name'";

			if ($id){
				$where .= " AND id<>$id";
			}

			$table = SYS_CONFIG_TBL;
			$status = $this->db->get_data(array(
				'where_clause'=> $where,
				'table_name'=> $table,
			));
			if (!$status){
				$this->error = $this->db->error;
			}
			return !$this->db->data;
		}

		function validateCategoryName($name, $id=-1){
			$table = SYS_CATEGORIES_TBL;
			return $this->validateName($name, $id, $table);
		}

		function validateConfigName($name, $id=-1){
			$table = SYS_CONFIG_TBL;
			return $this->validateName($name, $id, $table);
		}

		function idExists($id, $table){
			$data = array();
			$id = mysql_real_escape_string($id);

			$query = "SELECT * from $table WHERE id='$id'";

			if ($res=mysql_query($query)){
				while ($row = mysql_fetch_assoc($res)){
					$data[] = $row;
				}
				return $data[0];
			}
			else{
				$this->error[] = mysql_error();
				return false;
			}
		}

		function idConfigExists($id){
			$table = SYS_CONFIG_TBL;
			return $this->idExists($id, $table);
		}
		function idCatExists($id){
			$table = SYS_CATEGORIES_TBL;
			return $this->idExists($id, $table);
		}
	}
?>
