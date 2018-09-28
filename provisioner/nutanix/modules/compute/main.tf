
resource "nutanix_image" "CentOS" {
  name        = "CentOS"
  description = "Centos Cloud Image"
  source_uri  = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.qcow2"
}

resource "nutanix_virtual_machine" "bastion" {
  name = "bastion"
  description = "bastiotest1"
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 1024
  power_state          = "ON"
  nic_list = [
    {
      subnet_reference = {
        kind = "subnet"
        uuid = "33c97f77-2afd-4331-a2cc-542934091a48"
      }
      ip_endpoint_list = {
        ip   = "10.5.80.11"
        type = "ASSIGNED"
      }
    }
  ]
  disk_list = [{
     data_source_reference = [{
       kind = "image"      
       uuid = "${nutanix_image.CentOS.id}"
     }]
     device_properties = [{
       device_type = "DISK"
     }]
     disk_size_mib = 50000
   }]
  guest_customization_cloud_init_user_data = "${base64encode(file("cloud-init.conf"))}"
}

resource "nutanix_virtual_machine" "master" {
  name = "master"
  description = "master"
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 1024
  power_state          = "ON"
  nic_list = [
    {
      subnet_reference = {
        kind = "subnet"
        uuid = "33c97f77-2afd-4331-a2cc-542934091a48"
      }
      ip_endpoint_list = {
        ip   = "10.5.80.11"
        type = "ASSIGNED"
      }
    }
  ]
  disk_list = [{
     data_source_reference = [{
       kind = "image"      
       uuid = "${nutanix_image.CentOS.id}"
     }]
     device_properties = [{
       device_type = "DISK"
     }]
     disk_size_mib = 50000
   }]
  guest_customization_cloud_init_user_data = "${base64encode(file("cloud-init.conf"))}"
}

resource "nutanix_virtual_machine" "worker" {
  name = "worker"
  description = "worker"
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 1024
  power_state          = "ON"
  nic_list = [
    {
      subnet_reference = {
        kind = "subnet"
        uuid = "33c97f77-2afd-4331-a2cc-542934091a48"
      }
      ip_endpoint_list = {
        ip   = "10.5.80.11"
        type = "ASSIGNED"
      }
    }
  ]
  disk_list = [{
     data_source_reference = [{
       kind = "image"      
       uuid = "${nutanix_image.CentOS.id}"
     }]
     device_properties = [{
       device_type = "DISK"
     }]
     disk_size_mib = 50000
   }]
  guest_customization_cloud_init_user_data = "${base64encode(file("cloud-init.conf"))}"
}